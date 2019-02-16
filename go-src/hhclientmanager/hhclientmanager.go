package hhclientmanager

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/bitly/go-simplejson"
	uuid "github.com/satori/go.uuid"
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
}

func handleClientRequest(clientID *uuid.UUID, message *simplejson.Json) {
	//identify type of request
	requestType, err := message.Get("type").String()

	if err != nil {
		fmt.Println("Invalid request received.", err)
		return
	}

	switch requestType {
	case "NewPlayerRequest":
		handleNewPlayerRequest(clientID, message)
		break
	case "PositionUpdateRequest":
		handlePositionUpdateRequest(clientID, message)
		break
	default:
		fmt.Println("Unknown request type", requestType)
	}
}

func handleNewPlayerRequest(clientID *uuid.UUID, message *simplejson.Json) {
	message = message.Get("data")

	xBoardWidth := 1000.0
	yBoardWidth := 1000.0
	maxDirection := math.Pi * 2
	nickname, err := message.Get("nickname").String()

	if err != nil {
		fmt.Println(err)
		return
	}

	//TODO: check if the client has already joined
	//TODO: save xpos, ypos, dir and nickname to DB

	fmt.Println(nickname)

	x := rand.Float64() * xBoardWidth
	y := rand.Float64() * yBoardWidth
	direction := rand.Float64() * maxDirection

	response, err := createNewPlayerResponse(clientID, x, y, direction)

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

	fmt.Println("%lf %lf %lf", newX, newY, newDirection)
}
