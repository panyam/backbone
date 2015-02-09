package service

import (
	. "github.com/panyam/backbone/models"
)

type UserService struct {
	Cls IUserService
}

/**
* Get user info by ID
 */
func (s *UserService) GetUserById(id string) (*User, error) {
	return nil, nil
}

/**
* Get a user by username.
 */
func (s *UserService) GetUser(username string) (*User, error) {
	return nil, nil
}

/**
* Create a user with the given id and username.
 */
func (s *UserService) CreateUser(id string, username string) (*User, error) {
	return nil, nil
}
