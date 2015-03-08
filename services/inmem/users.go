package inmem

import (
	"errors"
	. "github.com/panyam/backbone/services/core"
)

type UserService struct {
	Cls         interface{}
	usersById   map[string]*User
	usersByName map[string]*User
}

func NewUserService() *UserService {
	svc := UserService{}
	svc.Cls = &svc
	svc.usersById = make(map[string]*User)
	svc.usersByName = make(map[string]*User)
	return &svc
}

/**
 * Get user info by ID
 */
func (s *UserService) GetUserById(id string) (*User, error) {
	return s.usersById[id], nil
}

/**
* Get a user by username.
 */
func (s *UserService) GetUser(username string) (*User, error) {
	return s.usersByName[username], nil
}

/**
* Create a user with the given id and username.
 */
func (s *UserService) CreateUser(id string, username string) (*User, error) {
	if _, ok := s.usersById[id]; ok {
		return nil, errors.New("User id already exists")
	}
	if _, ok := s.usersByName[username]; ok {
		return nil, errors.New("User name already exists")
	}
	newuser := NewUser(id, username)
	s.usersById[id] = newuser
	s.usersByName[username] = newuser
	return newuser, nil
}

func (s *UserService) SaveUser(user *User) error {
	return nil
}
