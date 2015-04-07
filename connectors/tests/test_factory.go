package connectors

import (
	"database/sql"
	_ "github.com/lib/pq"
	goclient "github.com/panyam/relay/clients/goclient"
	connector_core "github.com/panyam/relay/connectors/core"
	"github.com/panyam/relay/connectors/gorilla"
	authmw "github.com/panyam/relay/connectors/gorilla/middleware/auth"
	auth_core "github.com/panyam/relay/services/auth/core"
	auth_sqlds "github.com/panyam/relay/services/auth/sqlds"
	msg_core "github.com/panyam/relay/services/msg/core"
	msg_sqlds "github.com/panyam/relay/services/msg/sqlds"
	"log"
)

const factoryType = "sql"

func (s *TestSuite) CreateTestServer() connector_core.Server {
	server := gorilla.NewServer(s.ServerPort)

	s.DebugUserId = 666
	validator := authmw.NewDebugValidator(s.DebugUserId, s.serviceGroup.UserService)
	am := authmw.AuthMiddleware{Validators: []authmw.AuthValidator{validator}}
	server.SetAuthMiddleware(&am)
	return server
}

func (s *TestSuite) LoginClient() {
	s.client.EnableAuthentication(&goclient.DebugAuthenticator{Userid: 666})
}

func (s *TestSuite) CreateTestServices() (*msg_core.ServiceGroup, auth_core.IAuthService) {
	sg := msg_core.ServiceGroup{}
	var authService auth_core.IAuthService = nil
	if factoryType == "sql" {
		db, err := sql.Open("postgres", "user=test dbname=test sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		sg.TeamService = msg_sqlds.NewTeamService(db, &sg)
		sg.UserService = msg_sqlds.NewUserService(db, &sg)
		sg.ChannelService = msg_sqlds.NewChannelService(db, &sg)
		sg.MessageService = msg_sqlds.NewMessageService(db, &sg)

		authService = auth_sqlds.NewAuthService(db, sg.UserService, sg.TeamService)
	}
	return &sg, authService
}
