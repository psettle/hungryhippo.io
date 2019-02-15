package hhclientmanager

import (
	"fmt"

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
	default:
		fmt.Println("Unknown request type", requestType)
	}
}
