package core

import (
	"time"
)

type Team struct {
	Cls ITeam

	id string

	/**
	 * Name of the team - must be unique?
	 */
	name string

	created time.Time

	status int
}

/**
 * Users are actually unique message sources.  These are very generic and can be
 * from anywhere (eg a chat application, a github commit, a FB notification, an
 * email etc.
 */
type User struct {
	Cls IUser

	id string

	username string

	/**
	 * General purpose meta data.
	 */
	metaData map[string]interface{}

	/**
	 * List of teams to which this user belongs.
	 */
	teams []*Team

	addresses []IAddress
}

type Address struct {
	Cls    IAddress
	label  string
	domain string
	id     string
}

func NewUser() *User {
	user := User{}
	user.Cls = &user
	return &user
}

func (u *User) UserId() string {
	return u.id
}

func (u *User) Username() string {
	return u.username
}

func (u *User) SetUsername(value string) {
	u.username = value
}

func (u *User) Addresses() []IAddress {
	return u.addresses
}

func (u *User) AddAddress(address IAddress) {
}

func (u *User) GetMetaData(key string) interface{} {
	return nil
}

func (u *User) SetMetaData(key string, value interface{}) {
}

func NewAddress(domain string, id string, label string) *Address {
	address := Address{}
	address.Cls = &address
	address.domain = domain
	address.id = id
	address.label = label
	return &address
}

func (u *User) Self() IUser {
	return u.Cls
}

func (a *Address) Self() IAddress {
	return a.Cls
}

func (t *Team) Self() ITeam {
	return t.Cls
}

func (a *Address) Label() string {
	return a.label
}

func (a *Address) SetLabel(label string) {
	a.label = label
}

func (a *Address) Domain() string {
	return a.domain
}

func (a *Address) ID() string {
	return a.id
}

func (t *Team) TeamId() string {
	return t.id
}

func (t *Team) GetName() string {
	return t.name
}

func (t *Team) Created() time.Time {
	return t.created
}

func (t *Team) Status() int {
	return t.status
}
