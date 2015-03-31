package core

import (
	"github.com/panyam/relay/services/messaging/core"
)

type Server interface {
	Run()
	Stop()
	SetServiceGroup(sg *core.ServiceGroup)
}
