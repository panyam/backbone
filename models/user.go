package models

import (
	"time"
)

type Team struct {
	Id string

	/**
	 * Name of the team - must be unique?
	 */
	Name string

	Created time.Time

	Status int
}

/**
 * Users are actually unique message sources.  These are very generic and can be
 * from anywhere (eg a chat application, a github commit, a FB notification, an
 * email etc.
 */
type User struct {
	Id string

	Username string

	/**
	 * General purpose meta data.
	 */
	MetaData map[string]interface{}

	/**
	 * List of teams to which this user belongs.
	 */
	Teams []*Team

	Addresses []Address
}

type Address struct {
	Label  string
	Domain string
	Id     string
}

func NewUser(id string, username string) *User {
	user := User{Id: id, Username: username}
	return &user
}

func NewAddress(domain string, id string, label string) *Address {
	address := Address{}
	address.Domain = domain
	address.Id = id
	address.Label = label
	return &address
}
