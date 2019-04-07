package hhdatabase

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

type dbRecord struct {
	ip   string
	port string
}

type dbServers struct {
	dbs        []dbRecord
	lock       sync.Mutex
	primary    *pool.Pool
	electing   bool
	electLock  sync.Mutex
	looseCount int
	looseLock  sync.Mutex
}

var dbs dbServers

// initialize database with the line below NOTE: MUST HAVE DOCKER AND REDIS INSTALLED
//docker run -p 6379:6379 --name test-r -d redis
func init() {
	dbs = dbServers{}
	dbs.electing = false
	dbs.primary = nil
}

//NewDatabaseInstance accepts a new database instance to the system
func NewDatabaseInstance(ip string, port string) {
	dbInstance := dbRecord{}

	dbInstance.ip = ip
	dbInstance.port = port

	dbs.lock.Lock()
	fmt.Println("Database Register " + dbInstance.ip + ":" + dbInstance.port)
	//check if we already have that db
	for _, db := range dbs.dbs {
		if db.ip != dbInstance.ip {
			continue
		}

		if db.port != dbInstance.port {
			continue
		}

		//we already know about this db, perhaps the configurer added it twice
		return
	}

	dbs.dbs = append(dbs.dbs, dbInstance)

	//Run an election to ensure the new database is added to the network
	for {
		dbs.electLock.Lock()
		if dbs.electing {
			time.Sleep(10 * time.Millisecond)
		} else {
			dbs.electing = true
			fmt.Println("Add election starting" + dbInstance.ip + ":" + dbInstance.port)
			dbs.electLock.Unlock()
			break
		}
		dbs.electLock.Unlock()
	}
	electDatabase(nil)
	dbs.electLock.Lock()
	dbs.electing = false
	dbs.electLock.Unlock()
	fmt.Println("Database Add Complete " + dbInstance.ip + ":" + dbInstance.port)
	dbs.lock.Unlock()
}

//BeginOperation allocates a database connection for executing transactions
func BeginOperation() (*redis.Client, error) {

	dbs.electLock.Lock()
	if dbs.electing {
		dbs.electLock.Unlock()
		return nil, errors.New("Primary election in progress")
	}
	dbs.electLock.Unlock()

	dbs.lock.Lock()
	defer dbs.lock.Unlock()

	if dbs.primary == nil {
		return nil, errors.New("No primary database instance")
	}

	conn, err := dbs.primary.Get()

	if err == nil && conn != nil {
		dbs.looseLock.Lock()
		dbs.looseCount++
		dbs.looseLock.Unlock()
	}

	return conn, err
}

//Multi starts a transaction on the provided connection
func Multi(conn *redis.Client) error {
	return conn.Cmd("MULTI").Err
}

//Exec executes a transaction on the provided connection
//returns true if the transaction was applied, false if it failed due to watches
func Exec(conn *redis.Client) (bool, error) {
	response := conn.Cmd("EXEC")

	if response.Err != nil {
		return false, response.Err
	}

	/* Transaction didn't 'fail' but it might still have a null response
	   indicating that it wasn't applied due to conflict */
	if response.IsType(redis.Nil) {
		return false, nil
	}

	return true, nil
}

//EndOperation returns an allocated database connection
func EndOperation(conn *redis.Client) {
	dbs.primary.Put(conn)
	dbs.looseLock.Lock()
	if conn != nil {
		dbs.looseCount--
	}
	dbs.looseLock.Unlock()
}

//OnPrimaryFailure launches an election and blocks until the election is complete
func OnPrimaryFailure() {
	fmt.Println("OnPrimaryFailure")
	dbs.electLock.Lock()
	if dbs.electing {
		dbs.electLock.Unlock()
		fmt.Println("Election Already Running")
	} else {
		dbs.lock.Lock()
		dbs.electing = true
		fmt.Println("Started Election")
		dbs.electLock.Unlock()

		if len(dbs.dbs) > 0 && dbs.primary != nil {
			electDatabase(&dbs.dbs[0])
		} else {
			electDatabase(nil)
		}

		dbs.electLock.Lock()
		dbs.electing = false
		fmt.Println("Finished Election")
		dbs.electLock.Unlock()
		dbs.lock.Unlock()
	}
}

func electDatabase(deadInstance *dbRecord) {

	//start by safely closing the existing primary, if applicable
	if dbs.primary != nil {
		//we need to wait for all the allocated connections to be returned, so we don't leak any connections
		//this means all in-flight operations need to fail out
		//(new operations will auto-fail since the election flag is set)
		//wait a maximum of 100ms, more than that would indicate an error/leak of some kind
		for i := 0; i < 10; i++ {
			dbs.looseLock.Lock()
			fmt.Printf("Waiting for %d connections to be returned", dbs.looseCount)
			if dbs.looseCount == 0 {
				dbs.looseLock.Unlock()
				break
			} else {
				dbs.looseLock.Unlock()
				time.Sleep(10 * time.Millisecond)
			}

			if i == 9 {
				//timed out waiting for connections to be returned...
				fmt.Println("Failed to recover all connections in electDatabase")
			}
		}
		dbs.primary.Empty()
		dbs.primary = nil
	}

	if deadInstance != nil {
		//verify match, remove head
		if deadInstance.ip != dbs.dbs[0].ip ||
			deadInstance.port != dbs.dbs[0].port {
			//the dead instance is no longer the primary,
			//an election should have already taken place
			return
		}

		//else we need to delete the dead instance from the database queue
		dbs.dbs = append(dbs.dbs[:0], dbs.dbs[1:]...)
	}

	if len(dbs.dbs) == 0 {
		//no valid databases left :(
		//the next database to be added will be elected primary
		return
	}

	var err error
	dbs.primary, err = pool.New("tcp", dbs.dbs[0].ip+":"+dbs.dbs[0].port, 10)

	//if creating the pool failed, eject and start again
	if err != nil {
		electDatabase(&dbs.dbs[0])
		return
	}

	//restructure the network with the new primary:
	//note: all app servers will send duplicate commands
	//		the duplicates will be ignored by redis

	//convert primary to master
	//note: 'SLAVEOF' is deprecated, the 'modern' command is 'REPLICAOF'
	//we use 'SLAVEOF' because the windows binary doesn't support 'REPLICAOF' yet
	err = dbs.primary.Cmd("SLAVEOF", "NO", "ONE").Err
	fmt.Println("Configuring Primary: ", err)

	//tell all slaves who their primary is
	for i, db := range dbs.dbs {
		if i == 0 {
			//don't set primary as slave
			continue
		}

		conn, err := redis.Dial("tcp", db.ip+":"+db.port)
		if err != nil {
			//database is dead, will be cleared if it is ever up for primary
			fmt.Println(err)
			continue
		}

		err = conn.Cmd("SLAVEOF", dbs.dbs[0].ip, dbs.dbs[0].port).Err
		fmt.Println("Configuring Slave: ", err)
		if err != nil {
			fmt.Println(err)
		}

		conn.Close()
	}
}
