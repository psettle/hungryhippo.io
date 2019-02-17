package main

import (
	"encoding/json"
	"fmt"
	"github.com/mediocregopher/radix.v2/pool"
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

// Declare a global database variable to store the Redis connection pool.
var database *pool.Pool

// initialize database with the line below NOTE: MUST HAVE DOCKER AND REDIS INSTALLED
//docker run -p 6379:6379 --name test-r -d redis
func initDatabase() {
	var err error
	// Establish a pool of 10 connections to the Redis server listening on
	// port 6379 of the local machine.
	database, err = pool.New("tcp", "localhost:6379", 10)
	if err != nil {
		log.Panic(err)
	}
}

// Used for local testing purposes, will need to comment out when merged into project.
func main() {

	p1 := createPlayer("Chris")
	fmt.Println(p1)
	fmt.Println(p1.location())
	fmt.Println(p1.json())
	initDatabase()
	p1.savePlayer()
	getPlayer(81)
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

func (p *Player) savePlayer() {

	response := database.Cmd("HMSET",
		p.Id,
		"name", p.Name,
		"score", p.Score,
		"x", p.Location.X,
		"y", p.Location.Y,
		"direction", p.Location.Direction)

	if response.Err != nil {
		log.Fatal(response.Err)
	}
}

func getPlayer(id int) {
	reply, err := database.Cmd("HGETALL", id).Map()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}
