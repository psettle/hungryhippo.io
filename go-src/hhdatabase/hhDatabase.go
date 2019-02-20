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
	database, err = pool.New("tcp", "localhost:6379", 10)
	if err != nil {
		log.Panic(err)
	}
}

//BeginOperation allocates a database connection for executing transactions
func beginOperation() (*redis.Client, error) {
	return database.Get()
}

//EndOperation returns an allocated database connection
func endOperation(conn *redis.Client) {
	database.Put(conn)
}
