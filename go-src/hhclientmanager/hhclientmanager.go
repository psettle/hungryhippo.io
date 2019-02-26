package hhclientmanager

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bitly/go-simplejson"
	uuid "github.com/satori/go.uuid"
	"hungryhippo.io/go-src/hhdatabase"
	"hungryhippo.io/go-src/hhserver"
)

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

	//Create the player
	player := hhdatabase.CreatePlayer(clientID)
	player.Name = nickname
	player.Location.Centre.X = rand.Float64() * xBoardWidth
	player.Location.Centre.Y = rand.Float64() * yBoardWidth
	player.Location.Direction = rand.Float64() * maxDirection
	player.Score = 0

	var applied bool
	applied, err = createNewPlayer(player)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !applied {
		//if player was not created, do not return the new player response
		return
	}

	//successfully created the player, so respond now
	response, err := createNewPlayerResponse(player)
	if err != nil {
		fmt.Println(err)
		return
	}

	hhserver.SendJSON(clientID, response)
}

func handlePositionUpdateRequest(clientID *uuid.UUID, message *simplejson.Json) {
	//parse data
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

	//apply update
	_, err := updatePlayerPosition(hhdatabase.CreatePlayer(clientID), newX, newY, newDirection)
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

	fruit := hhdatabase.CreateFruit(&id)
	player := hhdatabase.CreatePlayer(clientID)

	//prepare a replacement fruit
	newFruitUUID := uuid.Must(uuid.NewV4())
	newFruit := hhdatabase.CreateFruit(&newFruitUUID)
	newFruit.Position.X = rand.Float64() * xBoardWidth
	newFruit.Position.Y = rand.Float64() * yBoardWidth

	//consume the fruit
	applied, err := consumeFruit(player, fruit, newFruit)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !applied {
		//opertion was invalid, don't proceed
		return
	}

	//the fruit was consumed, tell all the clients that it is gone.
	fruitConsumedMessage, respErr := createFruitConsumptionMessage(clientID, &id)
	if respErr != nil {
		fmt.Println(respErr)
		return
	}

	hhserver.SendJSONAll(fruitConsumedMessage)

	//and tell them about the new fruit
	newFruitMessage, mesgErr := createNewFruitMessage(newFruit)
	if mesgErr != nil {
		fmt.Println(mesgErr)
		return
	}

	hhserver.SendJSONAll(newFruitMessage)
}

func handleConsumePlayerRequest(clientID *uuid.UUID, message *simplejson.Json) {
	//parse out player ids
	consumerStr, consumerStrErr := message.Get("data").Get("consumer_id").String()
	consumerID, consumerErr := uuid.FromString(consumerStr)

	consumedStr, consumedStrErr := message.Get("data").Get("consumed_id").String()
	consumedID, consumedErr := uuid.FromString(consumedStr)

	switch {
	case consumerStrErr != nil:
		fmt.Println(consumerStrErr)
		return
	case consumerErr != nil:
		fmt.Println(consumerErr)
		return
	case consumedStrErr != nil:
		fmt.Println(consumedStrErr)
		return
	case consumedErr != nil:
		fmt.Println(consumedErr)
		return
	default:
		break
	}

	consumer := hhdatabase.CreatePlayer(&consumerID)
	consumed := hhdatabase.CreatePlayer(&consumedID)

	//apply consumption
	applied, err := consumePlayer(consumer, consumed)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !applied {
		//the operation was invalid
		return
	}

	//tell the consumer that they have grown
	playerConsumptionMessage, respErr := createConsumePlayerResponse(&consumer.ID, consumer.Score)
	if respErr != nil {
		fmt.Println(respErr)
		return
	}

	hhserver.SendJSON(&consumer.ID, playerConsumptionMessage)

	//tell the consumed that they died
	playerDeathMessage, mesgErr := createPlayerDeathMessage(&consumed.ID)
	if mesgErr != nil {
		fmt.Println(mesgErr)
		return
	}

	hhserver.SendJSON(&consumed.ID, playerDeathMessage)
}
