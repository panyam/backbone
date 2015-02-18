package gorilla

import (
	"github.com/gorilla/mux"
	"github.com/panyam/backbone/models"
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

type HttpHandlerFunc func(http.ResponseWriter, *http.Request)
type RequestHandlerFunc func(http.ResponseWriter, *http.Request, *RequestContext)

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

func (s *Server) DefaultMiddleware(requiresUser bool) *MiddlewareChain {
	return NewMiddlewareChain()
}

func (s *Server) Run() {
	r := mux.NewRouter()

	mwWithLogin := s.DefaultMiddleware(true)
	mwWithoutLogin := s.DefaultMiddleware(false)

	// /Users/sri/projects/go/src/github.com/panyam/backbone/clients
	http.Handle("/", r)
	r.HandleFunc("/", mwWithLogin.Apply(s.RootPageHandler()))
	r.HandleFunc("/login/", mwWithoutLogin.Apply(s.LoginPageHandler()))
	r.HandleFunc("/teams/", mwWithoutLogin.Apply(s.TeamListPageHandler()))
	r.HandleFunc("/teams/{id}", mwWithoutLogin.Apply(s.TeamPageHandler()))

	webHandler := http.StripPrefix("/webapp/", http.FileServer(http.Dir("../clients/web/app")))
	http.Handle("/webapp/", webHandler)

	apiRouter := s.createApiRouter(r)
	http.Handle("/api/", apiRouter)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func (s *Server) createApiRouter(parent *mux.Router) *mux.Router {
	apiRouter := parent.Path("/api").Subrouter()

	mwWithLogin := s.DefaultMiddleware(true)
	mwWithoutLogin := s.DefaultMiddleware(false)

	// Users/Login API
	accountRouter := apiRouter.PathPrefix("/account").Subrouter()
	accountRouter.HandleFunc("/register", mwWithoutLogin.Apply(s.AccountRegisterHandler()))
	accountRouter.HandleFunc("/login", mwWithoutLogin.Apply(s.AccountLoginHandler()))
	accountRouter.HandleFunc("/logout", mwWithoutLogin.Apply(s.AccountLogoutHandler()))

	// Teams API
	teamsRouter := apiRouter.Path("/teams/").Subrouter()
	teamsRouter.Methods("GET").HandlerFunc(mwWithLogin.Apply(s.GetTeamsHandler()))
	teamsRouter.Methods("POST").HandlerFunc(mwWithLogin.Apply(s.CreateTeamHandler()))

	teamRouter := apiRouter.PathPrefix("/teams/{id}").Subrouter()
	teamRouter.Methods("GET").HandlerFunc(mwWithLogin.Apply(s.TeamDetailsHandler()))
	teamRouter.Methods("PUT", "POST").HandlerFunc(mwWithLogin.Apply(s.UpdateTeamHandler()))
	teamRouter.Methods("DELETE").HandlerFunc(mwWithLogin.Apply(s.DeleteTeamHandler()))

	// Channels API
	channelsRouter := apiRouter.Path("/channels/").Subrouter()
	channelsRouter.Methods("GET").HandlerFunc(mwWithLogin.Apply(s.GetChannelsHandler()))
	channelsRouter.Methods("POST").HandlerFunc(mwWithLogin.Apply(s.CreateChannelHandler()))

	channelRouter := apiRouter.PathPrefix("/channels/{id}").Subrouter()
	channelRouter.Methods("GET").HandlerFunc(mwWithLogin.Apply(s.ChannelDetailsHandler()))
	channelRouter.Methods("PUT", "POST").HandlerFunc(mwWithLogin.Apply(s.UpdateChannelHandler()))
	channelRouter.Methods("DELETE").HandlerFunc(mwWithLogin.Apply(s.DeleteChannelHandler()))

	return apiRouter
}

func (s *Server) GetUser(request *http.Request) *models.User {
	request.ParseForm()
	log.Println("Form: ", request.Form)
	log.Println("Cookies: ", request.Cookies)
	return nil
}
