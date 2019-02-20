package hhclientmanager

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/bitly/go-simplejson"
	uuid "github.com/satori/go.uuid"
	"hungryhippo.io/go-src/hhdatabase"
	"hungryhippo.io/go-src/hhserver"
)

const xBoardWidth = 1000.0
const yBoardWidth = 1000.0
const maxDirection = math.Pi * 2

//HandleClients registers for websocket requests then accepts and responds to requests
func HandleClients() {
	//register for client requests
	requestHandler := make(chan interface{})

	hhserver.RegisterJSON(requestHandler)

	go func() {
		for {
			select {
			case requestI := <-requestHandler:
				request := requestI.(*hhserver.MessageJSON)
				go handleClientRequest(request.ClientID, request.Message)
				break
			}
		}
	}()

	//start position update task
	posUpdateTimer := time.NewTicker(time.Millisecond * 250)

	go func() {
		for {
			select {
			case <-posUpdateTimer.C:
				sendPositionUpdateMessage()
				break
			}
		}
	}()
}

func sendPositionUpdateMessage() {
	//generate the message
	message, err := createPositionUpdateMessage()

	if err != nil {
		//perhaps the database has crashed...
		fmt.Println(err)
		return
	}

	//send it to all clients
	hhserver.SendJSONAll(message)
}

func handleClientRequest(clientID *uuid.UUID, message *simplejson.Json) {
	//identify type of request
	requestType, err := message.Get("type").Int()

	if err != nil {
		fmt.Println("Invalid request received.", err)
		return
	}

	switch requestType {
	case newPlayerRequest:
		handleNewPlayerRequest(clientID, message)
		break
	case positionUpdateRequest:
		handlePositionUpdateRequest(clientID, message)
		break
	case consumeFruitRequest:
		handleConsumeFruitRequest(clientID, message)
		break
	case consumePlayerRequest:
		handleConsumePlayerRequest(clientID, message)
		break
	default:
		fmt.Println("Unknown request type", requestType)
	}
}

func handleNewPlayerRequest(clientID *uuid.UUID, message *simplejson.Json) {
	message = message.Get("data")

	nickname, err := message.Get("nickname").String()
	if err != nil {
		fmt.Println(err)
		return
	}

	player := hhdatabase.CreatePlayer(clientID)

	var exists bool
	exists, err = player.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	if exists {
		fmt.Println("Rejoin by existing player")
		return
	}

	player.Location.Centre.X = rand.Float64() * xBoardWidth
	player.Location.Centre.Y = rand.Float64() * yBoardWidth
	player.Location.Direction = rand.Float64() * maxDirection
	player.Name = nickname
	player.Score = 0

	//under the assumption that UUIDs are unique, there can be no conflict on save, so no need to retry or return a fail message
	//TODO: (technically there could be conflict if one client sent a new player request twice quickly)
	err = player.Save()
	if err != nil {
		fmt.Println(err)
		return
	}

	response, err := createNewPlayerResponse(player)
	if err != nil {
		fmt.Println(err)
		return
	}

	hhserver.SendJSON(clientID, response)
}

func handlePositionUpdateRequest(clientID *uuid.UUID, message *simplejson.Json) {
	location := message.Get("data").Get("location")

	newX, errX := location.Get("centre").Get("x").Float64()
	newY, errY := location.Get("centre").Get("y").Float64()
	newDirection, errDirection := location.Get("direction").Float64()

	switch {
	case errX != nil:
		fmt.Println(errX)
		return
	case errY != nil:
		fmt.Println(errY)
		return
	case errDirection != nil:
		fmt.Println(errDirection)
		return
	default:
		break
	}

	//load the associated player
	player := hhdatabase.CreatePlayer(clientID)
	err := player.Watch()
	if err != nil {
		fmt.Println(err)
		return
	}
	//need to defer a close call since watch player was called
	defer player.Close()

	exists, errExists := player.Load()
	if errExists != nil {
		fmt.Println(err)
		return
	}

	if !exists {
		fmt.Println("handlePositionUpdateRequest: Unknown player")
		return
	}

	//TODO: validate that movement was legal... (player didn't collide/get collided with)

	//save the new player position
	//TODO: this could fail if a collision happened between the .WatchPlayer call and the .Save call, the .Save call will fail in that case.
	player.Location.Centre.X = newX
	player.Location.Centre.Y = newY
	player.Location.Direction = newDirection
	err = player.Save()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func handleConsumeFruitRequest(clientID *uuid.UUID, message *simplejson.Json) {
	//which fruit is it?
	idstr, idstrErr := message.Get("data").Get("fruit_id").String()
	id, idErr := uuid.FromString(idstr)

	//check for parsing errors, indicates invalid id
	if idstrErr != nil {
		fmt.Println(idstrErr)
		return
	}

	if idErr != nil {
		fmt.Println(idErr)
		return
	}

	//delete that fruit
	fruit := hhdatabase.CreateFruit(&id)
	err := fruit.Watch()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fruit.Close()

	//try to delete it:
	var consumed bool
	consumed, err = fruit.Delete()
	if err != nil {
		fmt.Println(err)
		return
	}

	if !consumed {
		/* If this failed then the fruit has already been consumed, eat the message. */
		return
	}

	//the fruit was consumed, tell all the clients that it is gone.
	fruitConsumedMessage, respErr := createFruitConsumptionMessage(clientID, &id)
	if respErr != nil {
		fmt.Println(respErr)
		return
	}

	hhserver.SendJSONAll(fruitConsumedMessage)

	//generate a new fruit
	newFruitUUID := uuid.Must(uuid.NewV4())
	newFruit := hhdatabase.CreateFruit(&newFruitUUID)
	newFruit.Position.X = rand.Float64() * xBoardWidth
	newFruit.Position.Y = rand.Float64() * yBoardWidth
	err = newFruit.Save()

	if err != nil {
		log.Panic("Unique fruit could not be saved.")
	}

	newFruitMessage, mesgErr := createNewFruitMessage(newFruit)
	if mesgErr != nil {
		fmt.Println(mesgErr)
		return
	}

	hhserver.SendJSONAll(newFruitMessage)
}

func handleConsumePlayerRequest(clientID *uuid.UUID, message *simplejson.Json) {

}
