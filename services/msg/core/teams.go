package core

type CreateTeamRequest struct {
	Request

	Id int64

	Organization string

	Name string
}

type GetTeamsRequest struct {
	Request

	Organization string

	Offset int

	Count int
}

type GetTeamRequest struct {
	Request

	Id int64

	Organization string

	Name string
}

type DeleteTeamRequest struct {
	Request

	Team *Team
}

type TeamMembershipRequest struct {
	Request

	Team *Team

	Username string
}

type ITeamService interface {
	/**
	 * Removes all entries.
	 */
	RemoveAllTeams(request *Request)

	/**
	 * Create a team.
	 * If the ID is empty, then it is upto the backend to decide whether to
	 * throw an error or auto assign an ID.
	 * A valid Team object on return WILL have an ID if the backend can
	 * auto generate IDs
	 */
	CreateTeam(request *CreateTeamRequest) (*Team, error)

	/**
	 * Retrieve teams in a org
	 */
	GetTeamsInOrg(request *GetTeamsRequest) ([]*Team, error)

	/**
	 * Retrieve a team by Id
	 */
	GetTeamById(request *GetTeamRequest) (*Team, error)

	/**
	 * Retrieve a team by Name.
	 */
	GetTeamByName(request *GetTeamRequest) (*Team, error)

	/**
	 * Delete a team.
	 */
	DeleteTeam(request *DeleteTeamRequest) error

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
