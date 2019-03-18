package hhloadbalancer

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
)

type appServerRecord struct {
	ip   string
	port string
}

type appServers struct {
	servers []appServerRecord
	lock    sync.Mutex
}

var servers appServers

//StartServer starts the core http server for new clients;
func StartServer() {
	servers = appServers{}
	//setup static resource handlers
	htmlFileServer := http.FileServer(http.Dir("public/html-src/"))
	http.Handle("/", htmlFileServer)
	jsFileServer := http.FileServer(http.Dir("public/js-src/"))
	http.Handle("/js/", http.StripPrefix("/js/", jsFileServer))
	imgFileServer := http.FileServer(http.Dir("public/images/"))
	http.Handle("/img/", http.StripPrefix("/img/", imgFileServer))
	cssFileServer := http.FileServer(http.Dir("public/css-src/"))
	http.Handle("/css/", http.StripPrefix("/css/", cssFileServer))

	//setup dynamic resource handler
	http.HandleFunc("/ws", handleWebsocketRequest)
	http.HandleFunc("/as/", handleAppServerRegister)

	http.ListenAndServe(":80", nil)
}

func handleWebsocketRequest(w http.ResponseWriter, r *http.Request) {
	//someone is asking for an app server to open a websocket to:

	//first, check if the request has provided a dead server with it
	deadIP := r.Header.Get("dead-app-server-ip")
	deadPort := r.Header.Get("dead-app-server-port")
	for i, record := range servers.servers {
		if record.ip != deadIP {
			continue
		}

		if record.port != deadPort {
			continue
		}

		//found the dead server, delete it
		servers.lock.Lock()
		servers.servers = append(servers.servers[:i], servers.servers[i+1:]...)
		servers.lock.Unlock()
		break
	}

	//select a random app server
	servers.lock.Lock()
	if len(servers.servers) == 0 {
		servers.lock.Unlock()
		fmt.Println("No app servers")
		return
	}
	serverIndex := rand.Int() % len(servers.servers)
	appServer := servers.servers[serverIndex]
	response := []byte(`{
		"ip" : "` + appServer.ip + `",
		"port" : "` + appServer.port + `"
	}`)
	servers.lock.Unlock()

	//respond with ip and port to connect to
	w.Write(response)
}

func handleAppServerRegister(w http.ResponseWriter, r *http.Request) {
	appServer := appServerRecord{}

	values := r.URL.Query()

	appServer.ip = values.Get("app-server-ip")
	appServer.port = values.Get("app-server-port")

	fmt.Println("App Server Register " + appServer.ip + ":" + appServer.port)

	servers.lock.Lock()
	servers.servers = append(servers.servers, appServer)
	servers.lock.Unlock()

	//TODO: response with existing databases
}
