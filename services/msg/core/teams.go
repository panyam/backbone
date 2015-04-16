package core

type ITeamService interface {
	/**
	 * Removes all entries.
	 */
	RemoveAllTeams(request *Request)

	/**
	 * Create or update a team.
	 * If the ID is empty, then it is upto the backend to decide whether to
	 * throw an error or auto assign an ID.
	 * A valid Team object on return WILL have an ID if the backend can
	 * auto generate IDs
	 */
	SaveTeam(team *Team) (*Team, error)

	/**
	 * Retrieve teams in a org
	 */
	GetTeams(request *GetTeamsRequest) ([]*Team, error)

	/**
	 * Retrieve a team by either ID or Name and Organization
	 */
	GetTeam(team *Team) (*Team, error)

	/**
	 * Delete a team.
	 */
	DeleteTeam(team *Team) error

	/**
	 * Lets a user with the given username join a team (if allowed)
	 */
	JoinTeam(request *TeamMembershipRequest) (*User, error)

	/**
	 * Tells if a user belongs to a team.
	 */
	TeamContains(request *TeamMembershipRequest) bool

	/**
	 * Lets a user leave a team or be kicked out.
	 */
	LeaveTeam(request *TeamMembershipRequest) error
}

type GetTeamsRequest struct {
	Organization string

	Offset int

	Count int
}

type TeamMembershipRequest struct {
	Team *Team

	Username string
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

func NewTeam(id int64, org string, name string) *Team {
	team := Team{Organization: org, Name: name}
	team.Object = Object{Id: id}
	return &team
}

func NewTeamById(id int64) *Team {
	return NewTeam(id, "", "")
}
