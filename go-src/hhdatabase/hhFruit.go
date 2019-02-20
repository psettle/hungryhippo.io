package hhdatabase

import (
	"github.com/mediocregopher/radix.v2/redis"
	uuid "github.com/satori/go.uuid"
)

//Fruit defines a fruit score, id and position
type Fruit struct {
	ID       uuid.UUID
	Position Coordinate
	conn     *redis.Client
}

//CreateFruit initializes an empty fruit with an ID, it can then be loaded, saved, etc.
func CreateFruit(id *uuid.UUID) *Fruit {
	p := new(Fruit)
	p.ID = *id
	return p
}

//Watch puts a watch on a fruit entry, critical for load -> edit -> save operations,
//requires a close call after the transaction is complete.
func (p *Fruit) Watch() error {
	var err error

	if p.conn == nil {
		//there isn't an allocated connection yet... we need one.
		p.conn, err = beginOperation()

		if err != nil {
			return err
		}
	}

	reply := p.conn.Cmd("WATCH", p.ID.String())
	return reply.Err
}

//Load gets a single fruit from the database, returns true if the entry existed, false on not exists or error
func (p *Fruit) Load() (bool, error) {
	var conn interface {
		Cmd(cmd string, args ...interface{}) *redis.Resp
	}

	//load the fruit data...
	if p.conn != nil {
		//using allocated connection or...
		conn = p.conn
	} else {
		//using a standard conn...
		conn = database
	}

	reply, err := conn.Cmd("HMGET", "fruit_data", p.listKeys()).List()
	if err != nil {
		return false, err
	}

	//copy result to object
	err = p.fromList(reply)
	if err != nil {
		//error indicates the list was invalid, and since we know the command didn't fail, the fruit must not exist
		return false, nil
	}

	//no error
	return true, nil
}

//WatchFruits starts a watch on a set of fruit entries
//
// returns a client with which the watch started, can be used in LoadFruits to begin a bulk operation
// must be closed using endOperation(conn) to free up the allocated connection
//
//if ids == nil, a watch will be put on all fruits
func WatchFruits(ids []*uuid.UUID) (*redis.Client, error) {
	conn, err := beginOperation()
	if err != nil {
		return nil, err
	}

	//start with a watch on the fruit set
	reply := conn.Cmd("WATCH", "fruits")
	if reply.Err != nil {
		endOperation(conn)
		return nil, err
	}

	var idStrings []string
	if ids == nil {
		//ids was nil, load all members into a list
		idStrings, err = conn.Cmd("SMEMBERS", "fruits").List()
		if err != nil {
			endOperation(conn)
			return nil, err
		}
	} else {
		//else put all provided ids into a list
		for _, element := range ids {
			idStrings = append(idStrings, element.String())
		}
	}

	//put a watch on all members
	reply = conn.Cmd("WATCH", idStrings)
	if reply.Err != nil {
		endOperation(conn)
		return nil, err
	}

	return conn, nil
}

//UnWatchFruits ends a connection created by WatchFruits
func UnWatchFruits(conn *redis.Client) {
	endOperation(conn)
}

//LoadFruits returns
//a set of fruit objects for each id in ids
//a exists boolean for each id in ids
//a general error state for the operation
//
//if ids == nil, all fruits in the database will be returned
//if conn != nil, it will be used to execute load commands, this is to maintain a watch operation, the operation will be terminated on error
func LoadFruits(ids []*uuid.UUID, providedConn *redis.Client) ([]*Fruit, []bool, error) {
	var err error
	var conn *redis.Client
	//get a connection to load on
	if providedConn == nil {
		conn, err = beginOperation()
		if err != nil {
			return nil, nil, err
		}
		defer endOperation(conn)
	} else {
		conn = providedConn
	}

	//get ids to load
	var idStrings []string
	if ids == nil {
		//ids was nil, load all members into a list
		idStrings, err = conn.Cmd("SMEMBERS", "fruits").List()
		if err != nil {
			endOperation(conn)
			return nil, nil, err
		}
	} else {
		//else put all provided ids into a list
		for _, element := range ids {
			idStrings = append(idStrings, element.String())
		}
	}

	//fetch fruit data
	var fruits []*Fruit
	var exists []bool

	if ids == nil {
		//ids was nil, load all fruit data into a map, then let each fruit parse their own members
		entries, allErr := conn.Cmd("HGETALL", "fruit_data").Map()
		if allErr != nil {
			endOperation(conn)
			return nil, nil, allErr
		}

		for _, id := range idStrings {
			uuid, uuidErr := uuid.FromString(id)
			if uuidErr != nil {
				endOperation(conn)
				return nil, nil, uuidErr
			}

			p := CreateFruit(&uuid)
			p.conn = providedConn
			err = p.fromMap(entries)

			if err == nil {
				fruits = append(fruits, p)
				exists = append(exists, true)
			}
			//(Don't add invalid fruit entries)
		}
	} else {
		//ids was not nil, load all ids into a list, then let each fruit parse their own members
		var keys []string

		for _, id := range idStrings {
			uuid, uuidErr := uuid.FromString(id)
			if uuidErr != nil {
				endOperation(conn)
				return nil, nil, uuidErr
			}

			p := CreateFruit(&uuid)
			p.conn = providedConn
			fruits = append(fruits, p)

			keys = append(keys, p.listKeys()...)
		}

		list, allErr := conn.Cmd("HMGET", "fruit_data", keys).List()
		if allErr != nil {
			endOperation(conn)
			return nil, nil, allErr
		}

		loc := 0
		for _, p := range fruits {
			err = p.fromList(list[loc : loc+len(p.listKeys())])
			exists = append(exists, err == nil)

			loc += len(p.listKeys())
		}
	}

	return fruits, exists, nil
}

//Save saves a fruit entry to the database
func (p *Fruit) Save() error {
	//local copy of a connection to the database
	var conn *redis.Client

	if p.conn != nil {
		//use the already allocated conn
		conn = p.conn
	} else {
		//allocate a new conn for this operation
		var err error
		conn, err = beginOperation()
		if err != nil {
			return err
		}
		defer endOperation(conn)
	}

	//start a transaction
	response := conn.Cmd("MULTI")
	if response.Err != nil {
		return response.Err
	}

	id := p.ID.String()

	//put the fruit id into the set of fruit
	response = conn.Cmd("SADD", "fruits", id)
	if response.Err != nil {
		conn.Cmd("DISCARD")
		return response.Err
	}

	//save the fruit entry
	response = conn.Cmd("HMSET", "fruit_data", p.mapKeys())
	if response.Err != nil {
		conn.Cmd("DISCARD")
		return response.Err
	}

	//write to the fruit watch, to prevent concurrent access
	response = conn.Cmd("SET", id, "")
	if response.Err != nil {
		conn.Cmd("DISCARD")
		return response.Err
	}

	//commit the transaction
	response = conn.Cmd("EXEC")
	if response.Err != nil {
		return response.Err
	}

	//no errors
	return nil
}

//Delete removes a fruit entry from the database
func (p *Fruit) Delete() error {
	//local copy of a connection to the database
	var conn *redis.Client

	if p.conn != nil {
		//use the already allocated conn
		conn = p.conn
	} else {
		//allocate a new conn for this operation
		var err error
		conn, err = beginOperation()
		if err != nil {
			return err
		}
		defer endOperation(conn)
	}

	//start a transaction
	response := conn.Cmd("MULTI")
	if response.Err != nil {
		return response.Err
	}

	id := p.ID.String()

	//remove the fruit from the set of fruits
	response = conn.Cmd("SREM", "fruits", id)
	if response.Err != nil {
		conn.Cmd("DISCARD")
		return response.Err
	}

	//delete the fruit entry
	response = conn.Cmd("HREM", "fruit_data", p.listKeys())
	if response.Err != nil {
		conn.Cmd("DISCARD")
		return response.Err
	}

	//delete the fruit watch, to prevent concurrent access
	response = conn.Cmd("DEL", id)
	if response.Err != nil {
		conn.Cmd("DISCARD")
		return response.Err
	}

	//commit the transaction
	response = conn.Cmd("EXEC")
	if response.Err != nil {
		return response.Err
	}

	//no errors
	return nil
}

//Close frees resources associated with a fruit object
func (p *Fruit) Close() {
	//return the allocated connection if required
	if p.conn != nil {
		endOperation(p.conn)
	}
}

//Get database keys for elements
func (p *Fruit) listKeys() []string {
	id := p.ID.String()
	//TODO: correct list keys
	return []string{
		id + ":name",
		id + ":score",
		id + ":x",
		id + ":y",
		id + ":direction"}
}

func (p *Fruit) mapKeys() map[string]interface{} {
	//id := p.ID.String()

	return map[string]interface{}{}
	//TODO: correct return
	// id + ":name":      p.Name,
	// id + ":score":     p.Score,
	// id + ":x":         p.Location.Centre.X,
	// id + ":y":         p.Location.Centre.Y,
	// id + ":direction": p.Location.Direction}
}

func (p *Fruit) fromList(list []string) error {
	//TODO: correct from list
	// var err error

	// p.Name = list[0]
	// p.Score, err = strconv.Atoi(list[1])
	// if err != nil {
	// 	return err
	// }

	// p.Location.Centre.X, err = strconv.ParseFloat(list[2], 64)
	// if err != nil {
	// 	return err
	// }

	// p.Location.Centre.Y, err = strconv.ParseFloat(list[3], 64)
	// if err != nil {
	// 	return err
	// }

	// p.Location.Direction, err = strconv.ParseFloat(list[4], 64)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (p *Fruit) fromMap(entries map[string]string) error {
	//TODO: correct from map

	// var err error
	// id := p.ID.String()

	// p.Name = entries[id+":name"]
	// p.Score, err = strconv.Atoi(entries[id+":score"])
	// if err != nil {
	// 	return err
	// }
	// p.Location.Centre.X, err = strconv.ParseFloat(entries[id+":x"], 64)
	// if err != nil {
	// 	return err
	// }
	// p.Location.Centre.Y, err = strconv.ParseFloat(entries[id+":y"], 64)
	// if err != nil {
	// 	return err
	// }
	// p.Location.Direction, err = strconv.ParseFloat(entries[id+":direction"], 64)
	// if err != nil {
	// 	return err
	// }

	return nil
}
