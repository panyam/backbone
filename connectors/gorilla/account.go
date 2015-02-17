package gorilla

import (
	// "github.com/panyam/backbone/services"
	"log"
	"net/http"
)

func (s *Server) AccountRegisterHandler() HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		log.Println("Register account")
	}
}

func (s *Server) AccountLoginHandler() HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		log.Println("Login...")
	}
}

func (s *Server) AccountLogoutHandler() HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		log.Println("Logout ...")
	}
}
