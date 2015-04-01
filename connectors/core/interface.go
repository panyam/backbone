package core

import (
	authcore "github.com/panyam/relay/services/auth/core"
	msgcore "github.com/panyam/relay/services/msg/core"
)

type Server interface {
	Run()
	Stop()
	SetServiceGroup(sg *msgcore.ServiceGroup)
	SetAuthService(authSvc authcore.IAuthService)
}
