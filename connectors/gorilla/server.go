package gorilla

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/panyam/backbone/connectors/gorilla/middleware"
	"github.com/panyam/backbone/services/core"
	"log"
	"net/http"
)

type Server struct {
	serviceGroup *core.ServiceGroup
	store        *sessions.CookieStore
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) SetCookieStore(cs *sessions.CookieStore) {
	s.store = sessions.NewCookieStore([]byte("something-very-secret"))
}

func (s *Server) SetServiceGroup(sg *core.ServiceGroup) {
	s.serviceGroup = sg
}

func (s *Server) DefaultMiddleware(requiresUser bool) *middleware.MiddlewareChain {
	out := middleware.NewMiddlewareChain()
	out.AddRequestMiddleware(middleware.BodyParserMiddleware)
	// middleware to create jj
	if requiresUser {
	}
	return out
}

func (s *Server) Run() {
	log.Println("Starting http server...")
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

func (s *Server) Stop() {
	log.Println("Stopping http server...")
}

func (s *Server) createApiRouter(parent *mux.Router) *mux.Router {
	apiRouter := parent.Path("/api").Subrouter()

	mwWithLogin := s.DefaultMiddleware(true)
	mwWithoutLogin := s.DefaultMiddleware(false)

	// Users/Login API
	accountRouter := apiRouter.PathPrefix("/users").Subrouter()
	accountRouter.HandleFunc("/register/", mwWithoutLogin.Apply(s.AccountRegisterHandler()))
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

func (s *Server) GetUser(request *http.Request) *core.User {
	request.ParseForm()
	log.Println("Form: ", request.Form)
	log.Println("Cookies: ", request.Cookies)
	return nil
}
