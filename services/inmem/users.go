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
	svc.RemoveAllUsers()
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
func (s *UserService) GetUser(username string, team *Team) (*User, error) {
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

/**
 * Saves a user.
 * 	If the ID param is empty:
 * 		If username/team does not already exist a new one is created.
 * 		otherwise, it is updated and returned if override=true otherwise
 * 		false is returned.
 * 	Otherwise:
 * 		If username/team does not exist then it is written as is (Create or Update)
 * 		otherwise if IDs of curr and existing are different errow is thrown,
 * 		otherwise object is updated.
 */
func (s *UserService) SaveUser(user *User, override bool) error {
	if user.Id == "" {
	} else {
	}
	return nil
}

/**
 * Removes all entries.
 */
func (svc *UserService) RemoveAllUsers() {
	svc.usersById = make(map[string]*User)
	svc.usersByName = make(map[string]*User)
}
