package core

import (
	"time"
)

type Team struct {
	/**
	 * Unique system wide ID.
	 */
	Id string

	/**
	 * Name of this team.
	 */
	Name string

	/**
	 * Organization this team belongs to. (Org + Name must be unique)
	 */
	Organization string
}

/**
 * Users are actually unique message sources.  These are very generic and can be
 * from anywhere (eg a chat application, a github commit, a FB notification, an
 * email etc.
 */
type User struct {
	/**
	 * GUID.
	 */
	Id string

	/**
	 * The username that is unique within the team for this user.
	 */
	Username string

	/**
	 * The team the user belongs to.
	 * The combination of Team/Username *has* to be unique.
	 */
	Team *Team

	/**
	 * When the user was created.
	 */
	Created time.Time

	/**
	 * Status of the user account.
	 */
	Status int
}

func NewUser(id string, username string) *User {
	user := User{Id: id, Username: username}
	return &user
}

func NewTeam(id string, org string, name string) *Team {
	team := Team{Id: id, Organization: org, Name: name}
	return &team
}
