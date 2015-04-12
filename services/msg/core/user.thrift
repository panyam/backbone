
struct Team {
		i64 Id,
		string Name,
		string Organization
}

struct User {
		i64 Id,
		string username,
		Team team
}

exception NoSuchUser {
		i32 what
}

service UserService {
	/**
	 * Removes all entries.
	 */
	void RemoveAllUsers(),

	/**
	 * Get user info by ID
	 */
	User GetUserById(i64 id) throws (NoSuchUser error),

	/**
	 * Get a user by username in a particular team.
	 */
	User GetUser(String username, Team team) throws (NoSuchUser error),

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
	 error SaveUser(User user, bool override)
}
