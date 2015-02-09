package client

import (
	"github.com/panyam/backbone/models"
)

type Client interface {
	IsConnected() bool
	Connect() error
	Login(credentials map[string]string) error
	Logout() error
	Teams() ([]models.Team, error)
	CurrentUser() models.User
}
