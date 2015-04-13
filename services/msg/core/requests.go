package core

/**
 * Base struct for all the requests that will be handled by all the operations
 * that are handled and serviced by our service!
 */
type Request struct {
	// Every request must have a requestor
	Requestor *User
}

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
