package main

import (
	"fmt"
	"os"

	"github.com/psettle/hungryhippo.io/go-src/hhclientmanager"
	"github.com/psettle/hungryhippo.io/go-src/hhloadbalancer"
	"github.com/psettle/hungryhippo.io/go-src/hhserver"
)

const (
	appServer    = "appserver"
	loadBalancer = "loadbalancer"
)

const (
	programName = iota

	processTypeIndex = iota
	portIndex        = iota

	numArgs = iota
)

const buildMode = loadBalancer

func main() {
	if len(os.Args) < numArgs {
		fmt.Println("Need a process type and port to start")
		return
	}

	fmt.Println(os.Args[processTypeIndex])

	switch os.Args[processTypeIndex] {
	case appServer:
		startAppServer()
		break
	case loadBalancer:
		startLoadBalancer()
		break
	default:
		fmt.Println("Unknown process type")
	}

}

func startAppServer() {
	if len(os.Args) < numArgs {
		fmt.Println("Need a port to start an app server")
	}

	hhclientmanager.HandleClients()
	hhserver.StartServer(os.Args[portIndex])
}

func startLoadBalancer() {
	hhloadbalancer.StartServer()
}
