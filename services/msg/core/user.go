package core

type GetUserRequest struct {
	Request

	/**
	 * ID of the user to fetch.  If this is - then the username/team is looked
	 * up
	 */
	Id int64

	Username string

	Team *Team
}

type SaveUserRequest struct {
	Request

	/**
	 * The user to save.
	 */
	User *User

	/**
	 * Whether to override an existing user.
	 */
	Override bool
}

/**
 * Base service operations.  These dont care about authorization for now and
 * assume the user is authorized.  Authn (and possible Authz) have to be taken
 * care of seperately.
 */
type IUserService interface {
	/**
	 * Removes all entries.
	 */
	RemoveAllUsers(request *Request)

	/**
	 * Get user info by ID
	 */
	GetUserById(request *GetUserRequest) (*User, error)

	/**
	 * Get a user by username in a particular team.
	 */
	GetUser(request *GetUserRequest) (*User, error)

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
	SaveUser(request SaveUserRequest) error
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
