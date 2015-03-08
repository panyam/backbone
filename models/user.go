package models

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
	 * This system is unique system wide and is not settable by the client.
	 */
	Id string

	/**
	 * The user name for the user.
	 */
	Username string

	/**
	 * The team the user belongs to.
	 * The combination of Team/Username *has* to be unique.
	 * This gets set by the creator/owner of the team/org.  Why is the team a
	 * big deal?  We could have a "sri" in two different teams.
	 */
	Team *Team

	/**
	 * When the user was created.
	 */
	Created time.Time
}

func NewUser(id string, username string) *User {
	user := User{Id: id, Username: username}
	return &user
}

func NewTeam(id string, org string, name string) *Team {
	team := Team{Id: id, Organization: org, Name: name}
	return &team
}
