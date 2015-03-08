package services

import (
	"github.com/panyam/backbone/services/gae"
	"github.com/panyam/backbone/services/inmem"
)

var factoryType string = "gae"

func CreateServiceGroup() ServiceGroup {
	sg = ServiceGroup{}
	if factoryType == "inmem" {
		sg.ChannelService = inmem.NewChannelService()
		sg.UserService = inmem.NewUserService()
		sg.TeamService = inmem.NewTeamService()
		sg.MessageService = inmem.NewMessageService()
	} else if factoryType == "gae" {
		sg.ChannelService = gae.MockChannelService()
		sg.UserService = gae.MockUserService()
		sg.TeamService = gae.MockTeamService()
		sg.MessageService = gae.MockMessageService()
	}
	return sg
}
