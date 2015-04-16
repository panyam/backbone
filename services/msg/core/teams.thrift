
exception Error {
}

exception InvalidOperation {
		1: i32 what,
		2: string why
}

struct UTCTime {
		1: byte Year
		2: byte Month
		3: byte Day
		4: byte Hour
		5: byte Minute
		6: byte Second
		7: i32 MicroSecond
}

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


service TeamService {
	/**
	 * Removes all entries.
	 */
	void RemoveAllTeams(),

	/**
	 * Create or update a team.
	 * If the ID is empty, then it is upto the backend to decide whether to
	 * throw an error or auto assign an ID.
	 * A valid Team object on return WILL have an ID if the backend can
	 * auto generate IDs
	 */
	Team SaveTeam(1:Team team) throws (1:Error error),

	/**
	 * Retrieve teams in a org
	 */
	list<Team> GetTeams(1:GetTeamsRequest request) throws (1:Error error),

	/**
	 * Retrieve a team by either ID or Name and Organization
	 */
	Team GetTeam(1:Team team) throws (1:Error error),

	/**
	 * Delete a team.
	 */
	void DeleteTeam(1:Team team) throws (1:Error error),

	/**
	 * Lets a user with the given username join a team (if allowed)
	 */
	User JoinTeam(1:TeamMembershipRequest request) throws (1:Error error),

	/**
	 * Tells if a user belongs to a team.
	 */
	bool TeamContains(1: TeamMembershipRequest request),

	/**
	 * Lets a user leave a team or be kicked out.
	 */
	void LeaveTeam(1:TeamMembershipRequest request) throws (1:Error error)
}

struct Object {
	/**
	 * Tells if the object has been loaded from the entity store.
	 * This should be set by the entity manager when a load has finished.
	 */
	1:bool Loaded 

	/**
	 * Unique system wide ID.
	 */
	2: i64 Id 

	/**
	 * When the object was created.
	 */
	3:UTCTime Created 

	/**
	 * Status of the user account.
	 * 0 = valid and active
	 * everything else depends on the object type
	 */
	4:i32 Status
}

struct GetTeamsRequest {
	1:string Organization 
	2:i32 Offset
	3:i32 Count
}

struct TeamMembershipRequest {
	1: Team Team
	2: string Username
}

struct Team {
	1:Object Base 

	/**
	 * Name of this team.
	 */
	2:string Name

	/**
	 * Organization this team belongs to. (Org + Name must be unique)
	 */
	3: string Organization
}

