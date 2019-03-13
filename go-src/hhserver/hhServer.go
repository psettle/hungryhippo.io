package hhserver

import (
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
	imgFileServer := http.FileServer(http.Dir("public/images/"))
	http.Handle("/img/", http.StripPrefix("/img/", imgFileServer))
	cssFileServer := http.FileServer(http.Dir("public/css-src/"))
	http.Handle("/css/", http.StripPrefix("/css/", cssFileServer))

	//setup websocket entry point
	http.HandleFunc("/ws", socketRequestHandler)

	http.ListenAndServe(":80", nil)
}

//socketRequestHandler handles incoming requests to open a websocket from clients.
func socketRequestHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	//listen for messages
	go handleWebsocket(conn)
}
