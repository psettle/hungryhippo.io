package main

import (
	"encoding/json"
	"fmt"
	"github.com/mediocregopher/radix.v2/redis"
	"log"
	"math/rand"
	//"encoding/json"
)

type Player struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
	Location
}

type Location struct {
	X         int `json:"x"`
	Y         int `json:"y"`
	Direction int `json:"direction"`
}

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

	p1 := createPlayer("Chris")
	fmt.Println(p1)
	fmt.Println(p1.location())
	fmt.Println(p1.json())
	//resp := conn.Cmd("JSON.SET", "players", ".", p1.json())
	////Check the Err field of the *Resp object for any errors.
	//if resp.Err != nil {
	//	log.Fatal(resp.Err)
	//}
	//players := conn.Cmd("JSON.GET", "players")
	//var test Player
	//json.Unmarshal(players, &test)
	//fmt.Println(players)
}

// create player with randomly generated id and location
func createPlayer(name string) Player {
	player := Player{Id: generateID(), Name: name, Score: 0, Location: generateLocation()}
	return player
}

// Returns a random number between 0 and 999 to act as player Id, in future we will need to discuss how we plan to
// enforce that there are no duplicate id's.
func generateID() int {
	return rand.Intn(1000)
}

// generate a random location, x and y will be 0-999, will need to adjust this based on the size of our grid.
func generateLocation() Location {
	return Location{X: rand.Intn(1000), Y: rand.Intn(1000), Direction: 0}
}

// return string of the player's json representation
func (p *Player) json() string {
	player, err := json.Marshal(p)
	if err != nil {
		fmt.Println("player to json conversion failed:", err)
		return err.Error()
	} else {
		return string(player)
	}
}

// return string of player's location json representation
func (p *Player) location() string {
	location, err := json.Marshal(p.Location)
	if err != nil {
		fmt.Println("player location to json conversion failed:", err)
		return err.Error()
	} else {
		return string(location)
	}
}
