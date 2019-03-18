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
func StartServer(port string) {
	//register with the load balancer
	_, err := http.Get("http://localhost:80/as?app-server-ip=localhost&app-server-port=" + port)

	if err != nil {
		fmt.Println("Failed to register with app server: ", err)
		return
	}

	//setup static resource handlers
	htmlFileServer := http.FileServer(http.Dir("public/html-src/"))
	http.Handle("/", htmlFileServer)
	jsFileServer := http.FileServer(http.Dir("public/js-src/"))
	http.Handle("/js/", http.StripPrefix("/js/", jsFileServer))
	imgFileServer := http.FileServer(http.Dir("public/images/"))
	http.Handle("/img/", http.StripPrefix("/img/", imgFileServer))
	cssFileServer := http.FileServer(http.Dir("public/css-src/"))
	http.Handle("/css/", http.StripPrefix("/css/", cssFileServer))

	//we allow websocket requests from any url, this allows the load balancer to balance onto app servers
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	//setup websocket entry point
	http.HandleFunc("/ws", socketRequestHandler)

	http.ListenAndServe(":"+port, nil)
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
