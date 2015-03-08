package services

import (
	. "github.com/panyam/backbone/models"
	"github.com/panyam/backbone/services/gae"
	"github.com/panyam/backbone/services/inmem"
)

var factoryType string = "inmem"

func CreateChannelService() IChannelService {
	if factoryType == "inmem" {
		return inmem.NewChannelService()
	} else if factoryType == "gae" {
		return gae.NewChannelService()
	}
	return nil
}

func CreateMessageService() IMessageService {
	if factoryType == "inmem" {
		return inmem.NewMessageService()
	} else if factoryType == "gae" {
		return gae.NewMessageService()
	}
	return nil
}

func CreateUserService() IUserService {
	if factoryType == "inmem" {
		return inmem.NewUserService()
	} else if factoryType == "gae" {
		return gae.NewUserService()
	}
	return nil
}

func CreateTeamService() ITeamService {
	if factoryType == "inmem" {
		return inmem.NewTeamService()
	} else if factoryType == "gae" {
		return gae.NewTeamService()
	}
	return nil
}
