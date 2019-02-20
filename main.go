package main

import (
	"hungryhippo.io/go-src/hhclientmanager"
	"hungryhippo.io/go-src/hhserver"
)

func main() {
	hhclientmanager.HandleClients()
	hhserver.StartServer()
}
