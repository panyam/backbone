package gorilla

import (
	// "github.com/panyam/backbone/services"
	"log"
	"net/http"
)

func (s *Server) GetChannelsHandler() HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		log.Println("Get Channels")
	}
}

func (s *Server) CreateChannelHandler() HandlerFunc {
	return func(http.ResponseWriter, *http.Request) {
		log.Println("Create Channels")
	}
}

func (s *Server) ChannelDetailsHandler() HandlerFunc {
	return func(http.ResponseWriter, *http.Request) {
		log.Println("GetChannelDetails")
	}
}

func (s *Server) UpdateChannelHandler() HandlerFunc {
	return func(http.ResponseWriter, *http.Request) {
		log.Println("UpdateChannel")
	}
}

func (s *Server) DeleteChannelHandler() HandlerFunc {
	return func(http.ResponseWriter, *http.Request) {
		log.Println("DeleteChannel")
	}
}
