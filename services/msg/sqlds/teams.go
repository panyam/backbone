package sqlds

import (
	"database/sql"
	// "errors"
	"fmt"
	. "github.com/panyam/relay/services/msg/core"
	. "github.com/panyam/relay/utils"
)

type TeamService struct {
	Cls interface{}
	DB  *sql.DB
	SG  *ServiceGroup
}

const TEAMS_TABLE = "teams"
const TEAM_MEMBERS_TABLE = "team_members"

func NewTeamService(db *sql.DB, sg *ServiceGroup) *TeamService {
	svc := TeamService{}
	svc.Cls = &svc
	svc.SG = sg
	svc.DB = db
	svc.InitDB()
	return &svc
}

func (svc *TeamService) InitDB() {
	svc.SG.IDService.CreateDomain(&CreateDomainRequest{nil, "teamids", 1, 2})
	CreateTable(svc.DB, TEAMS_TABLE,
		[]string{
			"Id bigint PRIMARY KEY",
			"Organization TEXT NOT NULL",
			"Name TEXT NOT NULL",
		})

	CreateTable(svc.DB, TEAM_MEMBERS_TABLE,
		[]string{
			"UserId bigint NOT NULL REFERENCES users (Id) ON DELETE CASCADE",
			"TeamId bigint NOT NULL REFERENCES teams (Id) ON DELETE CASCADE",
			"Created TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp()",
			"Status INT DEFAULT (0)",
		},
		", CONSTRAINT unique_team_membership UNIQUE (UserId, TeamId)")
}

/**
 * Removes all entries.
 */
func (svc *TeamService) RemoveAllTeams(request *Request) {
	ClearTable(svc.DB, TEAMS_TABLE)
}

/**
 * Create a team.
 * If the ID is empty, then it is upto the backend to decide whether to
 * throw an error or auto assign an ID.
 * A valid Team object on return WILL have an ID if the backend can
 * auto generate IDs
 */
func (svc *TeamService) SaveTeam(team *Team) (*Team, error) {
	if team.Id == 0 {
		id2, err := svc.SG.IDService.NextID(&NextIDRequest{nil, "teamids"})
		if err != nil {
			return nil, err
		}
		team.Id = id2
	}
	query := fmt.Sprintf(`INSERT INTO %s ( Id, Organization, Name ) VALUES (%d, '%s', '%s')`, TEAMS_TABLE, team.Id, team.Organization, team.Name)
	_, err := svc.DB.Exec(query)
	if err != nil {
		return nil, err
	}
	return team, err
}

/**
 * Retrieve teams in a org
 */
func (svc *TeamService) GetTeams(request *GetTeamsRequest) ([]*Team, error) {
	return nil, nil
}

/**
 * Retrieve a team by ID.
 */
func (svc *TeamService) GetTeam(team *Team) (*Team, error) {
	var err error = nil
	if team.Id != 0 {
		query := fmt.Sprintf("SELECT Organization, Name from %s where Id = %d", TEAMS_TABLE, team.Id)
		row := svc.DB.QueryRow(query)

		err = row.Scan(&team.Organization, &team.Name)
	} else {
		query := fmt.Sprintf("SELECT Id from %s where Organization = '%s' and Name = '%s'",
			TEAMS_TABLE, team.Organization, team.Name)
		rows, err2 := svc.DB.Query(query)
		if err2 != nil {
			return nil, err2
		}
		defer rows.Close()
		rows.Next()
		err = rows.Scan(&team.Id)
		if err != nil {
			return nil, err
		}
	}
	if err == nil {
		team.Loaded = true
	}
	return team, err
}

/**
 * Delete a team.
 */
func (svc *TeamService) DeleteTeam(team *Team) error {
	return DeleteById(svc.DB, TEAMS_TABLE, team.Id)
}

/**
 * Lets a user to join a team (if allowed) and does not already exist.
 */
func (svc *TeamService) JoinTeam(request *TeamMembershipRequest) (*User, error) {
	return nil, nil
}

/**
 * Tells if a user belongs to a team.
 */
func (svc *TeamService) TeamContains(request *TeamMembershipRequest) bool {
	return false
}

/**
 * Lets a user leave a team or be kicked out.
 */
func (svc *TeamService) LeaveTeam(request *TeamMembershipRequest) error {
	return nil
}
