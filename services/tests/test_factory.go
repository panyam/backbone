package services

import (
	"appengine/aetest"
	"github.com/panyam/backbone/services/core"
	"github.com/panyam/backbone/services/gae"
	"github.com/panyam/backbone/services/inmem"
	"log"
)

var factoryType string = "gae"

func CreateServiceGroup() core.ServiceGroup {
	sg := core.ServiceGroup{}
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
	}
	return sg
}
