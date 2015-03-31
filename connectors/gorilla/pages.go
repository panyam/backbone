package gorilla

import (
	// "github.com/panyam/relay/services"
	. "github.com/panyam/relay/connectors/gorilla/common"
	"log"
	"net/http"
)

/**
 * Should see if we are logged in - then go to the team selection page otherwise
 * to login page.
 */
func (s *Server) RootPageHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		user := s.GetUser(request)
		log.Println("Root Page")
		if user == nil {
			// redirect to login page
			http.Redirect(rw, request, "/login/?next=/teams/", 302)
		} else {
			http.Redirect(rw, request, "/teams/", 302)
		}
	}
}

func (s *Server) LoginPageHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		request.ParseForm()
		log.Println("Login Page")
		log.Println("Form: ", request.Form)
	}
}

func (s *Server) TeamListPageHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		log.Println("Team List Page")
	}
}

func (s *Server) TeamPageHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		log.Println("Teams Page")
	}
}
