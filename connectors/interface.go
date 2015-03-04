package connectors

import (
	"github.com/panyam/backbone/services"
)

type Server interface {
	Run()
	Stop()
	SetTeamService(svc services.ITeamService)
	SetUserService(svc services.IUserService)
	SetChannelService(svc services.IChannelService)
	SetMessageService(svc services.IMessageService)
}
