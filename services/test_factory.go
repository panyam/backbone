package services

import (
	. "github.com/panyam/backbone/models"
	. "github.com/panyam/backbone/services/inmem"
)

func CreateChannelService() IChannelService {
	return NewChannelService()
}

func CreateMessageService() IMessageService {
	return NewMessageService()
}

func CreateUserService() IUserService {
	return NewUserService()
}

func CreateTeamService() ITeamService {
	return NewTeamService()
}
