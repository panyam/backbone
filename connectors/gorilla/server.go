package gorilla

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/panyam/relay/connectors/gorilla/common"
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

func (s *Server) MakeHandlerFunc(operation interface{},
	serviceRequestMaker common.ServiceRequestMaker,
	serviceResponsePresenter common.ServiceResponseMaker) common.HttpHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		service_request, err := serviceRequestMaker(request)
		if err != nil {
			serviceResponsePresenter(rw, request, nil, err)
		}

		service_method := operation.(func(interface{}) (interface{}, error))
		result, err := service_method(service_request)

		serviceResponsePresenter(rw, request, result, err)
	}
}

/**
 * Takes the response from a service operation (or an error) and presents it
 * back via the response mechanism.
 */
func (s *Server) SendServiceResponse(request *http.Request, rw http.ResponseWriter, result interface{}, err error) {
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

	// mwWithLogin := s.DefaultMiddleware(true)
	mwWithoutLogin := s.DefaultMiddleware(false)

	// Users/Login API
	accountRouter := apiRouter.PathPrefix("/users").Subrouter()
	accountRouter.HandleFunc("/registrations/{Id}/confirm/", mwWithoutLogin.Apply(s.AccountConfirmHandler()))
	accountRouter.HandleFunc("/register/", mwWithoutLogin.Apply(s.AccountRegisterHandler()))
	accountRouter.HandleFunc("/login", mwWithoutLogin.Apply(s.AccountLoginHandler()))
	accountRouter.HandleFunc("/logout", mwWithoutLogin.Apply(s.AccountLogoutHandler()))

	channelService := s.serviceGroup.ChannelService
	teamService := s.serviceGroup.TeamService
	// userService := s.serviceGroup.UserService
	// messageService := s.serviceGroup.MessageService

	// Teams API
	teamsRouter := apiRouter.Path("/teams/").Subrouter()
	teamsRouter.Methods("POST").HandlerFunc(s.MakeHandlerFunc(teamService.SaveTeam, SaveTeamRequestMaker, SaveTeamResponsePresenter))

	// Channel specific APi for a particular team
	teamChannelsRouter := apiRouter.PathPrefix("/teams/{teamId}/channels").Subrouter()
	teamChannelsRouter.Methods("GET").HandlerFunc(s.MakeHandlerFunc(channelService.GetChannels, GetChannelsRequestMaker, GetChannelsResponsePresenter))
	teamChannelsRouter.Methods("POST").HandlerFunc(s.MakeHandlerFunc(channelService.CreateChannel, CreateChannelRequestMaker, CreateChannelResponsePresenter))

	// Other teams API
	teamRouter := apiRouter.PathPrefix("/teams/{teamId}").Subrouter()
	teamRouter.Methods("GET").HandlerFunc(s.MakeHandlerFunc(teamService.GetTeam, GetTeamRequestMaker, GetTeamResponsePresenter))
	teamRouter.Methods("PUT", "POST").HandlerFunc(s.MakeHandlerFunc(teamService.SaveTeam, SaveTeamRequestMaker, SaveTeamResponsePresenter))
	teamRouter.Methods("DELETE").HandlerFunc(s.MakeHandlerFunc(teamService.DeleteTeam, DeleteTeamRequestMaker, DeleteTeamResponsePresenter))

	// Channels API
	channelsRouter := apiRouter.Path("/channels/").Subrouter()
	// channelsRouter.Methods("GET").HandlerFunc(s.MakeHandlerFunc(channelService.GetChannels, GetChannelsRequestMaker, GetChannelsResponsePresenter))
	channelsRouter.Methods("POST").HandlerFunc(s.MakeHandlerFunc(channelService.CreateChannel, CreateChannelRequestMaker, CreateChannelResponsePresenter))

	channelRouter := apiRouter.PathPrefix("/channels/{id}").Subrouter()
	channelRouter.Methods("GET").HandlerFunc(s.MakeHandlerFunc(channelService.GetChannelById, GetChannelRequestMaker, GetChannelResponsePresenter))
	channelRouter.Methods("PUT", "POST").HandlerFunc(s.MakeHandlerFunc(channelService.UpdateChannel, UpdateChannelRequestMaker, UpdateChannelResponsePresenter))
	channelRouter.Methods("DELETE").HandlerFunc(s.MakeHandlerFunc(channelService.DeleteChannel, DeleteChannelRequestMaker, DeleteChannelResponsePresenter))

	return apiRouter
}

func (s *Server) GetUser(request *http.Request) *msgcore.User {
	request.ParseForm()
	log.Println("Form: ", request.Form)
	log.Println("Cookies: ", request.Cookies)
	return nil
}
