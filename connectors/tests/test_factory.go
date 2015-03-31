package connectors

import (
	"database/sql"
	_ "github.com/lib/pq"
	connector_core "github.com/panyam/relay/connectors/core"
	"github.com/panyam/relay/connectors/gorilla"
	service_core "github.com/panyam/relay/services/core"
	"github.com/panyam/relay/services/sqlds"
	"log"
)

const factoryType = "sql"

func CreateTestServer() connector_core.Server {
	return gorilla.NewServer()
}

func CreateTestServiceGroup() *service_core.ServiceGroup {
	sg := service_core.ServiceGroup{}
	if factoryType == "sql" {
		db, err := sql.Open("postgres", "user=test dbname=test sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		sg.TeamService = sqlds.NewTeamService(db, &sg)
		sg.UserService = sqlds.NewUserService(db, &sg)
		sg.ChannelService = sqlds.NewChannelService(db, &sg)
		sg.MessageService = sqlds.NewMessageService(db, &sg)
	}
	return &sg
}
