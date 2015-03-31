package core

import (
	"time"
)

type Object struct {
	/**
	 * Unique system wide ID.
	 */
	Id int64

	/**
	 * When the object was created.
	 */
	Created time.Time

	/**
	 * Status of the user account.
	 * 0 = valid and active
	 * everything else = invalid
	 */
	Status int
}

type Team struct {
	Object

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
	Object

	/**
	 * The username that is unique within the team for this user.
	 */
	Username string

	/**
	 * The team the user belongs to.
	 * The combination of Team/Username *has* to be unique.
	 */
	Team *Team
}

func NewUser(id int64, username string, team *Team) *User {
	user := User{Username: username, Team: team}
	user.Object = Object{Id: id}
	return &user
}

func NewTeam(id int64, org string, name string) *Team {
	team := Team{Organization: org, Name: name}
	team.Object = Object{Id: id}
	return &team
}
