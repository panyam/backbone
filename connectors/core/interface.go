package core

import (
	"github.com/panyam/relay/services/msg/core"
)

type Server interface {
	Run()
	Stop()
	SetServiceGroup(sg *core.ServiceGroup)
}
