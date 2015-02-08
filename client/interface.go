package client

import (
	"github.com/panyam/backbone/core"
)

type Client interface {
	IsConnected() bool
	Connect() error
	Login(credentials map[string]string) error
	Logout() error
	Teams() ([]core.Team, error)
	CurrentUser() core.User
}
