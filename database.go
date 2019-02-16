package main

import (
	"fmt"
	"github.com/mediocregopher/radix.v2/redis"
	"log"
)

func main() {
	// Establish a connection to the Redis server listening on port 6379 of the
	// local machine. 6379 is the default port
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	// Importantly, use defer to ensure the connection is always properly
	// closed before exiting the main() function.
	defer conn.Close()

	// Send our command across the connection. The first parameter to Cmd()
	// is always the name of the Redis command (in this example HMSET),
	// optionally followed by any necessary arguments (in this example the
	// key, followed by the various hash fields and values).
	resp := conn.Cmd("HMSET", 1, "points", 0, "x", 0, "y", 0, "direction", 0)
	// Check the Err field of the *Resp object for any errors.
	if resp.Err != nil {
		log.Fatal(resp.Err)
	}

	fmt.Println("player 1 added")
}
