package client

import (
	"github.com/panyam/backbone/models"
)

type NativeClient struct {
	Cls         Client
	isConnected bool
	currentUser models.User
}

func (client *NativeClient) IsConnected() bool {
	return client.isConnected
}

func (client *NativeClient) Connect() error {
	client.isConnected = true
	return nil
}

func (client *NativeClient) Teams() ([]models.Team, error) {
	return nil, nil
}

func (client *NativeClient) CurrentUser() models.User {
	return client.currentUser
}

func (client *NativeClient) Login(credentials map[string]string) error {
	if credentials != nil {
		client.currentUser = models.NewUser()
		client.currentUser.SetUsername(credentials["username"])
	}
	return nil
}

func (client *NativeClient) Logout() error {
	client.currentUser = nil
	return nil
}
