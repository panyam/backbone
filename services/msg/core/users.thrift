
include "common.thrift"
include "teams.thrift"

/**
 * Base service operations.  These dont care about authorization for now and
 * assume the user is authorized.  Authn (and possible Authz) have to be taken
 * care of seperately.
 */
service UserService {
	/**
	 * Removes all entries.
	 */
	void RemoveAllUsers()

	/**
	 * Get user info by ID or username/team
	 */
	User GetUser(1:User user) throws (1:Error error)

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
	void SaveUser(1:SaveUserRequest request) throws (1:Error error)
}

/**
 * Users are actually unique message sources.  These are very generic and can be
 * from anywhere (eg a chat application, a github commit, a FB notification, an
 * email etc.
 */
struct User {
	1: Object Base

	/**
	 * The username that is unique within the team for this user.
	 */
	2: string Username 

	/**
	 * The team the user belongs to.
	 * The combination of Team/Username *has* to be unique.
	 */
	3: Team Team
}

struct SaveUserRequest {
	/**
	 * The user to save.
	 */
	1: User User

	/**
	 * Whether to override an existing user.
	 */
	2: bool Override
}

