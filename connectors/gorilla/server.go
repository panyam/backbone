package gorilla

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/panyam/relay/connectors/gorilla/middleware"
	authmw "github.com/panyam/relay/connectors/gorilla/middleware/auth"
	authcore "github.com/panyam/relay/services/auth/core"
	msgcore "github.com/panyam/relay/services/msg/core"

	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Server struct {
	serviceGroup   *msgcore.ServiceGroup
	authService    authcore.IAuthService
	store          *sessions.CookieStore
	authMiddleware *authmw.AuthMiddleware
	DebugUserId    int64
	Port           int
}

func NewServer(port int) *Server {
	rand.Seed(time.Now().UTC().UnixNano())
	return &Server{Port: port}
}

func (s *Server) SetCookieStore(cs *sessions.CookieStore) {
	s.store = sessions.NewCookieStore([]byte("something-very-secret"))
}

func (s *Server) SetServiceGroup(sg *msgcore.ServiceGroup) {
	s.serviceGroup = sg
}

func (s *Server) SetAuthService(authSvc authcore.IAuthService) {
	s.authService = authSvc
}

func (s *Server) SetAuthMiddleware(am *authmw.AuthMiddleware) {
	s.authMiddleware = am
}

func (s *Server) DefaultMiddleware(requiresUser bool) *middleware.MiddlewareChain {
	out := middleware.NewMiddlewareChain()
	out.AddRequestMiddleware(middleware.BodyParserMiddleware)
	if requiresUser {
		out.AddMiddleware(s.authMiddleware)
	}
	/*
		out.AddResponseMiddleware(func(rw http.ResponseWriter, request *http.Request, context common.IRequestContext) error {
			errs := context.Errors()
			log.Println("Errors: ", errs)
			if len(errs) == 0 {
				output := context.Get("output").(interface{})
				utils.SendJsonResponse(rw, output)
			} else {
				http.Error(rw, errs[0].Error(), context.Get("StatusCode").(int))
			}
			return nil
		})
	*/
	return out
}

func (s *Server) Run() {
	log.Println("Starting http server...")
	r := mux.NewRouter()

	mwWithLogin := s.DefaultMiddleware(true)
	mwWithoutLogin := s.DefaultMiddleware(false)

	// /Users/sri/projects/go/src/github.com/panyam/relay/clients
	http.Handle("/", r)
	r.HandleFunc("/", mwWithLogin.Apply(s.RootPageHandler()))
	r.HandleFunc("/login/", mwWithoutLogin.Apply(s.LoginPageHandler()))
	r.HandleFunc("/teams/", mwWithoutLogin.Apply(s.TeamListPageHandler()))
	r.HandleFunc("/teams/{id}", mwWithoutLogin.Apply(s.TeamPageHandler()))

	webHandler := http.StripPrefix("/webapp/", http.FileServer(http.Dir("../clients/web/app")))
	http.Handle("/webapp/", webHandler)

	apiRouter := s.createApiRouter(r)
	http.Handle("/api/", apiRouter)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil))
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
	accountRouter.HandleFunc("/registrations/{Id}/confirm/", mwWithoutLogin.Apply(s.AccountConfirmHandler()))
	accountRouter.HandleFunc("/register/", mwWithoutLogin.Apply(s.AccountRegisterHandler()))
	accountRouter.HandleFunc("/login", mwWithoutLogin.Apply(s.AccountLoginHandler()))
	accountRouter.HandleFunc("/logout", mwWithoutLogin.Apply(s.AccountLogoutHandler()))

	// Teams API
	teamsRouter := apiRouter.Path("/teams/").Subrouter()
	teamsRouter.Methods("GET").HandlerFunc(mwWithLogin.Apply(s.GetTeamsHandler()))
	teamsRouter.Methods("POST").HandlerFunc(mwWithLogin.Apply(s.CreateTeamHandler()))

	// Channel specific APi for a particular team
	// teamChannelsRouter := apiRouter.PathPrefix("/teams/{teamid}/channels").Subrouter()
	// teamChannelsRouter.Methods("POST").HandlerFunc(mwWithLogin.Apply(s.CreateChannelHandler()))

	// Other teams API
	teamRouter := apiRouter.PathPrefix("/teams/{teamid}").Subrouter()
	teamRouter.Methods("GET").HandlerFunc(mwWithLogin.Apply(s.TeamDetailsHandler()))
	teamRouter.Methods("PUT", "POST").HandlerFunc(mwWithLogin.Apply(s.SaveTeamHandler()))
	teamRouter.Methods("DELETE").HandlerFunc(mwWithLogin.Apply(s.DeleteTeamHandler()))

	// Channels API
	channelsRouter := apiRouter.Path("/channels/").Subrouter()
	// channelsRouter.Methods("GET").HandlerFunc(mwWithLogin.Apply(s.GetChannelsHandler()))
	channelsRouter.Methods("POST").HandlerFunc(mwWithLogin.Apply(s.CreateChannelHandler()))

	channelRouter := apiRouter.PathPrefix("/channels/{id}").Subrouter()
	channelRouter.Methods("GET").HandlerFunc(mwWithLogin.Apply(s.ChannelDetailsHandler()))
	channelRouter.Methods("PUT", "POST").HandlerFunc(mwWithLogin.Apply(s.UpdateChannelHandler()))
	channelRouter.Methods("DELETE").HandlerFunc(mwWithLogin.Apply(s.DeleteChannelHandler()))

	return apiRouter
}

func (s *Server) GetUser(request *http.Request) *msgcore.User {
	request.ParseForm()
	log.Println("Form: ", request.Form)
	log.Println("Cookies: ", request.Cookies)
	return nil
}
