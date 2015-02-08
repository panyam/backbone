package core

import (
	"time"
)

type SimpleTeam struct {
	Cls Team

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
type SimpleUser struct {
	Cls User

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

	addresses []Address
}

type SimpleAddress struct {
	Cls    Address
	label  string
	domain string
	id     string
}

func NewUser() *SimpleUser {
	user := SimpleUser{}
	user.Cls = &user
	return &user
}

func (u *SimpleUser) UserId() string {
	return u.id
}

func (u *SimpleUser) Username() string {
	return u.username
}

func (u *SimpleUser) SetUsername(value string) {
	u.username = value
}

func (u *SimpleUser) Addresses() []Address {
	return u.addresses
}

func (u *SimpleUser) AddAddress(address *Address) {
}

func (u *SimpleUser) GetMetaData(key string) interface{} {
	return nil
}

func (u *SimpleUser) SetMetaData(key string, value interface{}) {
}

func NewAddress(domain string, id string, label string) *SimpleAddress {
	address := SimpleAddress{}
	address.Cls = &address
	address.domain = domain
	address.id = id
	address.label = label
	return &address
}

func (u *SimpleUser) Self() User {
	return u.Cls
}

func (a *SimpleAddress) Self() Address {
	return a.Cls
}

func (t *SimpleTeam) Self() Team {
	return t.Cls
}

func (a *SimpleAddress) Label() string {
	return a.label
}

func (a *SimpleAddress) SetLabel(label string) {
	a.label = label
}

func (a *SimpleAddress) Domain() string {
	return a.domain
}

func (a *SimpleAddress) ID() string {
	return a.id
}

func (t *SimpleTeam) TeamId() string {
	return t.id
}

func (t *SimpleTeam) GetName() string {
	return t.name
}

func (t *SimpleTeam) Created() time.Time {
	return t.created
}

func (t *SimpleTeam) Status() int {
	return t.status
}
