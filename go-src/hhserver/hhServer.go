package hhserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//StartServer starts the core http server for new clients;
func StartServer() {
	//setup static resource handlers
	htmlFileServer := http.FileServer(http.Dir("public/html-src/"))
	http.Handle("/", htmlFileServer)
	jsFileServer := http.FileServer(http.Dir("public/js-src/"))
	http.Handle("/js/", http.StripPrefix("/js/", jsFileServer))

	//setup websocket entry point
	http.HandleFunc("/ws", SocketRequestHandler)

	http.ListenAndServe(":80", nil)
}

//SocketRequestHandler handles incoming requests to open a websocket from clients.
func SocketRequestHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	conn.WriteJSON("hello")

	//listen for messages
	go listenWebsocket(conn)
}

type msg struct {
	Message string `json:"message"`
}

//Read thread for a websocket connection
func listenWebsocket(conn *websocket.Conn) {
	for {
		m := msg{}

		err := conn.ReadJSON(&m)
		if err != nil {
			fmt.Println("Error reading json.", err)
		}

		fmt.Printf("Got message: %s\n", m.Message)

		m.Message = "Howdy!"

		if err = conn.WriteJSON(m); err != nil {
			fmt.Println(err)
		}
	}
}
