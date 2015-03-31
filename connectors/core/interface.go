package core

import (
	"github.com/panyam/relay/services/core"
)

type Server interface {
	Run()
	Stop()
	SetServiceGroup(sg *core.ServiceGroup)
}
