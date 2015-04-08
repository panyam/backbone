package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/panyam/relay/connectors/gorilla"
	authmw "github.com/panyam/relay/connectors/gorilla/middleware/auth"
	auth_sqlds "github.com/panyam/relay/services/auth/sqlds"
	msgcore "github.com/panyam/relay/services/msg/core"
	msgsqlds "github.com/panyam/relay/services/msg/sqlds"
	"log"
)

func CreateServer() *gorilla.Server {
	db, err := sql.Open("postgres", "user=test dbname=test sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	sg := msgcore.ServiceGroup{}
	sg.TeamService = msgsqlds.NewTeamService(db, &sg)
	sg.UserService = msgsqlds.NewUserService(db, &sg)
	sg.ChannelService = msgsqlds.NewChannelService(db, &sg)
	sg.MessageService = msgsqlds.NewMessageService(db, &sg)

	authService := auth_sqlds.NewAuthService(db, sg.UserService, sg.TeamService)

	server := gorilla.NewServer(3000)
	server.SetServiceGroup(&sg)
	server.SetAuthService(authService)

	// some dummy data for testing
	debugUserId := int64(666)
	validator := authmw.NewDebugValidator(debugUserId, sg.UserService)
	am := authmw.AuthMiddleware{Validators: []authmw.AuthValidator{validator}}
	server.SetAuthMiddleware(&am)

	testTeam, err := sg.TeamService.CreateTeam(1, "testorg", "testteam")
	if err != nil {
		testTeam, err = sg.TeamService.GetTeamById(1)
	}
	log.Println("TestTeam: ", testTeam, err)
	debugUser := msgcore.NewUser(debugUserId, "debuguser", testTeam)
	sg.UserService.SaveUser(debugUser, false)
	return server
}

func main() {
	server := CreateServer()
	server.Run()
}
