package gorilla

import (
	// "github.com/panyam/backbone/services"
	"log"
	"net/http"
)

func (s *Server) GetChannelsHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		log.Println("Get Channels")
	}
}

func (s *Server) CreateChannelHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		log.Println("Create Channels")
	}
}

func (s *Server) ChannelDetailsHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		log.Println("GetChannelDetails")
	}
}

func (s *Server) UpdateChannelHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		log.Println("UpdateChannel")
	}
}

func (s *Server) DeleteChannelHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		log.Println("DeleteChannel")
	}
}
