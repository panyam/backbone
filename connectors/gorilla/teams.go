package gorilla

import (
	// "github.com/panyam/relay/services"
	. "github.com/panyam/relay/connectors/gorilla/common"
	"log"
	"net/http"
)

func (s *Server) CreateTeamHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		if context.Get("user") == nil {
		}
		log.Println("Create Teams")
	}
}

func (s *Server) GetTeamsHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		log.Println("Get Teams")
	}
}

func (s *Server) TeamDetailsHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		log.Println("GetTeamDetails")
	}
}

func (s *Server) UpdateTeamHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		log.Println("UpdateTeam")
	}
}

func (s *Server) DeleteTeamHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		log.Println("DeleteTeam")
	}
}
