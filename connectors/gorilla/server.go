package gorilla

import (
	"github.com/gorilla/mux"
	"github.com/panyam/backbone/services"
	"log"
	"net/http"
)

type Server struct {
	teamService    services.ITeamService
	userService    services.IUserService
	channelService services.IChannelService
	messageService services.IMessageService
}

type HandlerFunc func(http.ResponseWriter, *http.Request)

func NewServer() *Server {
	return &Server{}
}

func (s *Server) SetTeamService(svc services.ITeamService) {
	s.teamService = svc
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

func (s *Server) createApiRouter(parent *mux.Router) *mux.Router {
	apiRouter := parent.Path("/api").Subrouter()

	// Users/Login API
	accountRouter := apiRouter.PathPrefix("/account").Subrouter()
	accountRouter.HandleFunc("/register", s.AccountRegisterHandler())
	accountRouter.HandleFunc("/login", s.AccountLoginHandler())
	accountRouter.HandleFunc("/logout", s.AccountLogoutHandler())

	// Teams API
	teamsRouter := apiRouter.Path("/teams/").Subrouter()
	teamsRouter.Methods("GET").HandlerFunc(s.GetTeamsHandler())
	teamsRouter.Methods("POST").HandlerFunc(s.CreateTeamHandler())

	teamRouter := apiRouter.PathPrefix("/teams/{id}").Subrouter()
	teamRouter.Methods("GET").HandlerFunc(s.TeamDetailsHandler())
	teamRouter.Methods("PUT", "POST").HandlerFunc(s.UpdateTeamHandler())
	teamRouter.Methods("DELETE").HandlerFunc(s.DeleteTeamHandler())

	// Channels API
	channelsRouter := apiRouter.Path("/channels/").Subrouter()
	channelsRouter.Methods("GET").HandlerFunc(s.GetChannelsHandler())
	channelsRouter.Methods("POST").HandlerFunc(s.CreateChannelHandler())

	channelRouter := apiRouter.PathPrefix("/channels/{id}").Subrouter()
	channelRouter.Methods("GET").HandlerFunc(s.ChannelDetailsHandler())
	channelRouter.Methods("PUT", "POST").HandlerFunc(s.UpdateChannelHandler())
	channelRouter.Methods("DELETE").HandlerFunc(s.DeleteChannelHandler())

	return apiRouter
}
