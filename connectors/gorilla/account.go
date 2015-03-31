package gorilla

import (
	// "github.com/panyam/relay/services"
	. "github.com/panyam/relay/connectors/gorilla/common"
	"log"
	"net/http"
)

func (s *Server) AccountRegisterHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		log.Println("Register account")
	}
}

func (s *Server) AccountLoginHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		log.Println("Login...")
	}
}

func (s *Server) AccountLogoutHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		log.Println("Logout ...")
	}
}
