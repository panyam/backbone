package gorilla

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func StaticHandler(path string) func(http.ResponseWriter, *http.Request) {
	return nil
}

func (s *Server) Run() {
	r := mux.NewRouter()

	// /Users/sri/projects/go/src/github.com/panyam/backbone/clients
	http.Handle("/", r)
	http.Handle("/webapp/", http.StripPrefix("/webapp/", http.FileServer(http.Dir("../clients/web/app"))))
	log.Fatal(http.ListenAndServe(":3000", nil))
}
