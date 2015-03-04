package main

import (
	// "github.com/panyam/backbone/connectors/gocraft"
	"github.com/panyam/backbone/connectors"
	"github.com/panyam/backbone/connectors/gorilla"
)

func CreateServer() connectors.Server {
	return gorilla.NewServer()
}

func main() {
	server := CreateServer()
	server.Run()
}
