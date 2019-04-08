package hhloadbalancer

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type appServerRecord struct {
	ip   string
	port string
}

type appServers struct {
	servers []appServerRecord
	lock    sync.Mutex
}

type dbRecord struct {
	ip   string
	port string
}

type dbServers struct {
	dbs  []dbRecord
	lock sync.Mutex
}

var servers appServers
var dbs dbServers

//StartServer starts the core http server for new clients;
func StartServer() {
	servers = appServers{}
	dbs = dbServers{}
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
	http.HandleFunc("/db/", handleDatabaseRegister)

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
		fmt.Println("App Server Dead " + deadIP + ":" + deadPort)
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

	//check if we already have that server
	for _, server := range servers.servers {
		if server.ip != appServer.ip {
			continue
		}

		if server.port != appServer.port {
			continue
		}

		//we already know about this server, perhaps it was close and restarted while no clients existed.
		return
	}

	go distributeExistingDatabases(&appServer)
	servers.servers = append(servers.servers, appServer)
	servers.lock.Unlock()

	//TODO: response with existing databases
}

func handleDatabaseRegister(w http.ResponseWriter, r *http.Request) {
	dbInstance := dbRecord{}

	values := r.URL.Query()

	dbInstance.ip = values.Get("db-ip")
	dbInstance.port = values.Get("db-port")

	fmt.Println("Database Register " + dbInstance.ip + ":" + dbInstance.port)

	dbs.lock.Lock()

	//check if we already have that db
	for _, db := range dbs.dbs {
		if db.ip != dbInstance.ip {
			continue
		}

		if db.port != dbInstance.port {
			continue
		}

		//we already know about this db, perhaps the configurer added it twice
		return
	}

	distributeNewDatabase(&dbInstance)
	dbs.dbs = append(dbs.dbs, dbInstance)
	dbs.lock.Unlock()
}

func distributeExistingDatabases(appserver *appServerRecord) {
	//take a sec to ensure the app server has initialized
	time.Sleep(100 * time.Millisecond)
	base := "http://" + appserver.ip + ":" + appserver.port + "/db/?"

	dbs.lock.Lock()
	for _, db := range dbs.dbs {
		query := base + "db-ip=" + db.ip + "&db-port=" + db.port
		_, err := http.Get(query)

		if err != nil {
			//err typically means that the app server has gone down
			//this is not the time to recover from that
			//we will recover when the re-routed clients tell the load balancer that
			//the app server died
		}
	}
	dbs.lock.Unlock()
}

func distributeNewDatabase(db *dbRecord) {
	base := "db-ip=" + db.ip + "&db-port=" + db.port

	for _, appServer := range servers.servers {

		query := "http://" + appServer.ip + ":" + appServer.port + "/db/?" + base

		_, err := http.Get(query)

		if err != nil {
			//err typically means that the app server has gone down
			//this is not the time to recover from that
			//we will recover when the re-routed clients tell the load balancer that
			//the app server died
		}
	}
}
