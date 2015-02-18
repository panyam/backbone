package gorilla

import (
	// "github.com/panyam/backbone/services"
	"log"
	"net/http"
)

func (s *Server) CreateTeamHandler() HandlerFunc {
	return func(http.ResponseWriter, *http.Request) {
		s.teamService.Create
		log.Println("Create Teams")
	}
}

func (s *Server) GetTeamsHandler() HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		log.Println("Get Teams")
	}
}

func (s *Server) TeamDetailsHandler() HandlerFunc {
	return func(http.ResponseWriter, *http.Request) {
		log.Println("GetTeamDetails")
	}
}

func (s *Server) UpdateTeamHandler() HandlerFunc {
	return func(http.ResponseWriter, *http.Request) {
		log.Println("UpdateTeam")
	}
}

func (s *Server) DeleteTeamHandler() HandlerFunc {
	return func(http.ResponseWriter, *http.Request) {
		log.Println("DeleteTeam")
	}
}
