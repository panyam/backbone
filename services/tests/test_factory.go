package services

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/panyam/relay/services/core"
	// "appengine/aetest"
	// "github.com/panyam/relay/services/gae"
	// "github.com/panyam/relay/services/inmem"
	"github.com/panyam/relay/services/sqlds"
	"log"
)

var factoryType string = "sql"

func CreateServiceGroup() core.ServiceGroup {
	sg := core.ServiceGroup{}
	/*
		if factoryType == "inmem" {
			sg.ChannelService = inmem.NewChannelService()
			sg.UserService = inmem.NewUserService()
			sg.TeamService = inmem.NewTeamService()
			sg.MessageService = inmem.NewMessageService()
		} else if factoryType == "gae" {
			ctx, err := aetest.NewContext(nil)
			if err != nil {
				log.Println("NewContext error: ", err)
			}
			sg.ChannelService = gae.NewChannelService(ctx)
			sg.UserService = gae.NewUserService(ctx)
			sg.TeamService = gae.NewTeamService(ctx)
			sg.MessageService = gae.NewMessageService(ctx)
		} else */
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
	return sg
}
