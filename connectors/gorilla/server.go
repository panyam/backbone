package gorilla

import (
	"github.com/gorilla/mux"
	"github.com/panyam/backbone/services"
	"log"
	"net/http"
)

type Server struct {
	userService    services.IUserService
	channelService services.IChannelService
	messageService services.IMessageService
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) SetUserService(svc services.IUserService) {
	s.userService = svc
}

func (s *Server) SetChannelService(svc services.IChannelService) {
	s.channelService = svc
}

func (s *Server) SetMessageService(svc services.IMessageService) {
	s.messageService = svc
}

type HandlerFunc func(http.ResponseWriter, *http.Request)

func (s *Server) GetChannelsHandler() HandlerFunc {
	return func(http.ResponseWriter, *http.Request) {
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

func (s *Server) createApiRouter(parent *mux.Router) *mux.Router {
	apiRouter := parent.Path("/api").Subrouter()
	channelsRouter := apiRouter.Path("/channels/").Subrouter()
	channelsRouter.Methods("GET").HandlerFunc(s.GetChannelsHandler())
	channelsRouter.Methods("POST").HandlerFunc(s.CreateChannelHandler())

	channelRouter := apiRouter.PathPrefix("/posts/{id}").Subrouter()
	channelRouter.Methods("GET").HandlerFunc(s.ChannelDetailsHandler())
	channelRouter.Methods("PUT", "POST").HandlerFunc(s.UpdateChannelHandler())
	channelRouter.Methods("DELETE").HandlerFunc(s.DeleteChannelHandler())

	return apiRouter
}

func (s *Server) Run() {
	r := mux.NewRouter()

	// /Users/sri/projects/go/src/github.com/panyam/backbone/clients
	http.Handle("/", r)

	webHandler := http.StripPrefix("/webapp/", http.FileServer(http.Dir("../clients/web/app")))
	http.Handle("/webapp/", webHandler)

	apiRouter := s.createApiRouter(r)
	http.Handle("/api/", apiRouter)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
