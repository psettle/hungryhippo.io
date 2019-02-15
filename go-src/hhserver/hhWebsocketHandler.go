package hhserver

import (
	"fmt"

	"github.com/dustin/go-broadcast"
	uuid "github.com/satori/go.uuid"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
)

//type for recording the low level connection data about clients
type activeConnection struct {
	clientID *uuid.UUID
	conn     *websocket.Conn
}

type MessageJSON struct {
	ClientID *uuid.UUID
	Message  *simplejson.Json
}

type activeConnectionSet struct {
	connections    map[uuid.UUID]activeConnection
	newConn        chan *activeConnection
	sendMessage    chan *MessageJSON
	sendAllMessage chan *simplejson.Json
	receiveMessage chan *MessageJSON
	closeConn      chan *uuid.UUID
	rxListeners    broadcast.Broadcaster
}

var activeConnections activeConnectionSet

func init() {
	//start the activeConnections manager
	globalWebsocketSync()
}

//RegisterJSON adds a listener for incoming events from all clients
func RegisterJSON(listener chan interface{}) {
	activeConnections.rxListeners.Register(listener)
}

//SendJSON sends a message to a client asynchronously
func SendJSON(clientID *uuid.UUID, request *simplejson.Json) {
	go sendJSONBlocking(clientID, request)
}

//SendJSONAll sends a message to all clients asynchronously
func SendJSONAll(request *simplejson.Json) {
	go sendAllJSONBlocking(request)
}

//main handler for websocket events, synced to one function so activeConnections in not accessed concurrently
func globalWebsocketSync() {
	//list of active clients
	activeConnections = activeConnectionSet{}
	activeConnections.connections = make(map[uuid.UUID]activeConnection)
	activeConnections.newConn = make(chan *activeConnection)
	activeConnections.sendMessage = make(chan *MessageJSON)
	activeConnections.sendAllMessage = make(chan *simplejson.Json)
	activeConnections.receiveMessage = make(chan *MessageJSON)
	activeConnections.closeConn = make(chan *uuid.UUID)
	activeConnections.rxListeners = broadcast.NewBroadcaster(100)

	//handle all client events, synced here for safety on activeConnections
	go func() {
		for {
			select {
			case conn := <-activeConnections.newConn:
				//add the new connection
				activeConnections.connections[*conn.clientID] = *conn
				break
			case toSend := <-activeConnections.sendMessage:
				//send to the connection if it still exists
				if val, ok := activeConnections.connections[*toSend.ClientID]; ok {
					val.conn.WriteJSON(toSend.Message)
				}
				break
			case toSendAll := <-activeConnections.sendAllMessage:
				//send to each connection
				for k := range activeConnections.connections {
					activeConnections.connections[k].conn.WriteJSON(toSendAll)
				}

				break
			case received := <-activeConnections.receiveMessage:
				//process the message
				activeConnections.rxListeners.Submit(received)

				//check if there is another message
				if val, ok := activeConnections.connections[*received.ClientID]; ok {
					go websocketJSONReceive(&val)
				}

				break
			case toClose := <-activeConnections.closeConn:
				//remove the connection
				delete(activeConnections.connections, *toClose)
				break
			}
		}
	}()
}

//handle a new websocket connection
func handleWebsocket(conn *websocket.Conn) {
	//create a record for the connection
	connectionRecord := activeConnection{}
	connectionRecord.conn = conn

	clientID := uuid.Must(uuid.NewV4())
	connectionRecord.clientID = &clientID
	activeConnections.newConn <- &connectionRecord

	//start a receive task
	go websocketJSONReceive(&connectionRecord)
}

//receive a single JSON message and pass to main handler thread
func websocketJSONReceive(conn *activeConnection) {
	request := simplejson.New()

	err := conn.conn.ReadJSON(request)

	if err != nil {
		fmt.Println("Error reading json.", err)

		//something is wrong with the connection, kill the handlers
		activeConnections.closeConn <- conn.clientID
		return
	}

	activeConnections.receiveMessage <- &MessageJSON{conn.clientID, request}
}

//Send a JSON message to all current clients, blocks until complete
func sendAllJSONBlocking(request *simplejson.Json) {
	activeConnections.sendAllMessage <- request
}

//Send a JSON message to one client, blocks until complete
func sendJSONBlocking(clientID *uuid.UUID, request *simplejson.Json) {
	/* Send the request over the channel */
	message := MessageJSON{}
	message.ClientID = clientID
	message.Message = request

	activeConnections.sendMessage <- &message
}
