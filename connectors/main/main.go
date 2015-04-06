package main

import (
	// "github.com/panyam/relay/connectors/gocraft"
	"github.com/panyam/relay/connectors"
	"github.com/panyam/relay/connectors/gorilla"
)

func CreateServer() connectors.Server {
	gorilla := gorilla.NewServer(3000)
	gorilla.SetServiceGroup(sg)
	db, err := sql.Open("postgres", "user=test dbname=test sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	sg.TeamService = sqlds.NewTeamService(db, &sg)
	sg.UserService = sqlds.NewUserService(db, &sg)
	sg.ChannelService = sqlds.NewChannelService(db, &sg)
	sg.MessageService = sqlds.NewMessageService(db, &sg)
}

func main() {
	server := CreateServer()
	server.Run()
}
