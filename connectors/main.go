package main

import (
	// "github.com/panyam/backbone/connectors/gocraft"
	"github.com/panyam/backbone/connectors/gorilla"
	"github.com/panyam/backbone/services"
)

type Server interface {
	Run()
	SetUserService(svc services.IUserService)
	SetChannelService(svc services.IChannelService)
	SetMessageService(svc services.IMessageService)
}

func CreateServer() Server {
	return gorilla.NewServer()
}

func main() {
	server := CreateServer()
	server.Run()
}
