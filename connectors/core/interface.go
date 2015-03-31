package core

import (
	"github.com/panyam/backbone/services/core"
)

type Server interface {
	Run()
	Stop()
	SetServiceGroup(sg *core.ServiceGroup)
}
