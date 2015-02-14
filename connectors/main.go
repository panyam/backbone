package main

import (
	// "github.com/panyam/backbone/connectors/gocraft"
	"github.com/panyam/backbone/connectors/gorilla"
)

func main() {
	server := gorilla.NewServer()
	server.Run()
}
