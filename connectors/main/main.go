package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/panyam/relay/connectors/gorilla"
	authmw "github.com/panyam/relay/connectors/gorilla/middleware/auth"
	auth_sqlds "github.com/panyam/relay/services/auth/sqlds"
	msg_core "github.com/panyam/relay/services/msg/core"
	msg_sqlds "github.com/panyam/relay/services/msg/sqlds"
	"log"
)

func CreateServer() *gorilla.Server {
	db, err := sql.Open("postgres", "user=test dbname=test sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	sg := msg_core.ServiceGroup{}
	sg.TeamService = msg_sqlds.NewTeamService(db, &sg)
	sg.UserService = msg_sqlds.NewUserService(db, &sg)
	sg.ChannelService = msg_sqlds.NewChannelService(db, &sg)
	sg.MessageService = msg_sqlds.NewMessageService(db, &sg)

	authService := auth_sqlds.NewAuthService(db, sg.UserService, sg.TeamService)

	server := gorilla.NewServer(3000)
	server.SetServiceGroup(&sg)
	server.SetAuthService(authService)

	server.DebugUserId = 666
	validator := authmw.NewDebugValidator(server.DebugUserId, sg.UserService)
	am := authmw.AuthMiddleware{Validators: []authmw.AuthValidator{validator}}
	server.SetAuthMiddleware(&am)

	return server
}

func main() {
	server := CreateServer()
	server.Run()
}
