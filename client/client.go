package client

import (
	"github.com/panyam/backbone/core"
)

type NativeClient struct {
	Cls         Client
	isConnected bool
	currentUser core.User
}

func (client *NativeClient) IsConnected() bool {
	return client.isConnected
}

func (client *NativeClient) Connect() error {
	client.isConnected = true
	return nil
}

func (client *NativeClient) Teams() ([]core.Team, error) {
	return nil, nil
}

func (client *NativeClient) CurrentUser() core.User {
	return client.currentUser
}

func (client *NativeClient) Login(credentials map[string]string) error {
	if credentials != nil {
		client.currentUser = core.NewUser()
		client.currentUser.SetUsername(credentials["username"])
	}
	return nil
}

func (client *NativeClient) Logout() error {
	client.currentUser = nil
	return nil
}
