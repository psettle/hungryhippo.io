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

	//start game update task
	gameUpdateTimer := time.NewTicker(time.Millisecond * 100)

	go func() {
		for {
			select {
			case <-gameUpdateTimer.C:
				sendGamestateUpdateMessage()
				break
			}
		}
	}()
}

func sendGamestateUpdateMessage() {
	//generate the message
	message, err := createGamestateUpdateMessage()

	if err != nil {
		fmt.Println(err)
		return
	}

	//send it to all clients
	hhserver.SendJSONAll(message)
}

func handleClientRequest(clientID *uuid.UUID, message *simplejson.Json) {
	//a nil message indicates an abrupt client disconnect
	if message == nil {
		err := handleClientDisconnect(clientID)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	//identify type of request
	requestType, err := message.Get("type").Int()

	if err != nil {
		fmt.Println("Invalid request received.", err)
		return
	}

	switch requestType {
	case newPlayerRequest:
		err = handleNewPlayerRequest(clientID, message)
		break
	case positionUpdateRequest:
		err = handlePositionUpdateRequest(clientID, message)
		break
	case consumeFruitRequest:
		err = handleConsumeFruitRequest(clientID, message)
		break
	case consumePlayerRequest:
		err = handleConsumePlayerRequest(clientID, message)
		break
	default:
		fmt.Println("Unknown request type", requestType)
	}

	if err != nil {
		fmt.Println(err)
		return
	}
}

func handleClientDisconnect(clientID *uuid.UUID) error {
	player := hhdatabase.CreatePlayer(clientID)
	_, err := deletePlayer(player)
	return err
}

func handleNewPlayerRequest(clientID *uuid.UUID, message *simplejson.Json) error {
	message = message.Get("data")

	nickname, err := message.Get("nickname").String()
	if err != nil {
		return err
	}

	//Create the player
	player := hhdatabase.CreatePlayer(clientID)
	player.Name = nickname
	player.Location.Centre.X = rand.Float64() * xBoardWidth
	player.Location.Centre.Y = rand.Float64() * yBoardWidth
	player.Location.Direction = rand.Float64() * maxDirection
	player.Score = 0

	//create the new fruit
	fruitID, err := uuid.NewV4()
	if err != nil {
		return err
	}

	fruit := hhdatabase.CreateFruit(&fruitID)
	fruit.Position.X = rand.Float64() * xBoardWidth
	fruit.Position.Y = rand.Float64() * xBoardWidth

	var applied bool
	applied, err = createNewPlayer(player, fruit)
	if err != nil {
		return err
	}

	if !applied {
		//if player was not created, do not return the new player response
		return nil
	}

	//successfully created the player, so respond now
	response, err := createNewPlayerResponse(player)
	if err != nil {
		return err
	}

	hhserver.SendJSON(clientID, response)
	return nil
}

func handlePositionUpdateRequest(clientID *uuid.UUID, message *simplejson.Json) error {
	//parse data
	location := message.Get("data").Get("location")
	newX, errX := location.Get("centre").Get("x").Float64()
	newY, errY := location.Get("centre").Get("y").Float64()
	newDirection, errDirection := location.Get("direction").Float64()

	switch {
	case errX != nil:
		return errX
	case errY != nil:
		return errY
	case errDirection != nil:
		return errDirection
	default:
		break
	}

	//apply update
	_, err := updatePlayerPosition(hhdatabase.CreatePlayer(clientID), newX, newY, newDirection)
	return err
}

func handleConsumeFruitRequest(clientID *uuid.UUID, message *simplejson.Json) error {
	//which fruit is it?
	idstr, idstrErr := message.Get("data").Get("fruit_id").String()
	id, idErr := uuid.FromString(idstr)

	//check for parsing errors, indicates invalid id
	if idstrErr != nil {
		return idstrErr
	}

	if idErr != nil {
		return idErr
	}

	fruit := hhdatabase.CreateFruit(&id)
	player := hhdatabase.CreatePlayer(clientID)

	//prepare a replacement fruit
	newFruitUUID := uuid.Must(uuid.NewV4())
	newFruit := hhdatabase.CreateFruit(&newFruitUUID)
	newFruit.Position.X = rand.Float64() * xBoardWidth
	newFruit.Position.Y = rand.Float64() * yBoardWidth

	//consume the fruit
	_, err := consumeFruit(player, fruit, newFruit)
	return err
}

func handleConsumePlayerRequest(clientID *uuid.UUID, message *simplejson.Json) error {
	//parse out player ids
	consumerStr, consumerStrErr := message.Get("data").Get("consumer_id").String()
	consumerID, consumerErr := uuid.FromString(consumerStr)

	consumedStr, consumedStrErr := message.Get("data").Get("consumed_id").String()
	consumedID, consumedErr := uuid.FromString(consumedStr)

	switch {
	case consumerStrErr != nil:
		return consumerStrErr
	case consumerErr != nil:
		return consumerErr
	case consumedStrErr != nil:
		return consumedStrErr
	case consumedErr != nil:
		return consumedErr
	default:
		break
	}

	consumer := hhdatabase.CreatePlayer(&consumerID)
	consumed := hhdatabase.CreatePlayer(&consumedID)

	//apply consumption
	_, err := consumePlayer(consumer, consumed)
	return err
}
