package hhdatabase

import (
	"log"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

// Declare a global database variable to store the Redis connection pool.
var database *pool.Pool

// initialize database with the line below NOTE: MUST HAVE DOCKER AND REDIS INSTALLED
//docker run -p 6379:6379 --name test-r -d redis
func init() {
	var err error
	// Establish a pool of 10 connections to the Redis server listening on
	// port 6379 of the local machine.
	database, err = pool.New("tcp", "redis:6379", 10)
	if err != nil {
		log.Panic(err)
	}
}

//BeginOperation allocates a database connection for executing transactions
func BeginOperation() (*redis.Client, error) {
	return database.Get()
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
	database.Put(conn)
}
