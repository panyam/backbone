package sqlds

import (
	"database/sql"
	"fmt"
	. "github.com/panyam/relay/services/msg/core"
	. "github.com/panyam/relay/utils"
	"log"
)

type UserService struct {
	Cls interface{}
	SG  *ServiceGroup
	DB  *sql.DB
}

const USERS_TABLE = "users"

func NewUserService(db *sql.DB, sg *ServiceGroup) *UserService {
	svc := UserService{}
	svc.SG = sg
	svc.Cls = &svc
	svc.DB = db
	svc.InitDB()
	return &svc
}

func (svc *UserService) InitDB() {
	CreateTable(svc.DB, USERS_TABLE,
		[]string{
			"Id bigint PRIMARY KEY",
			"Username TEXT NOT NULL",
			"TeamId bigint NOT NULL REFERENCES teams (Id) ON DELETE CASCADE",
			"Created TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp()",
			"Status INT DEFAULT (0)",
		},
		", CONSTRAINT unique_username_team UNIQUE (Username, TeamId)")
}

/**
 * Removes all entries.
 */
func (svc *UserService) RemoveAllUsers() {
	ClearTable(svc.DB, USERS_TABLE)
}

/**
 * Get user info by ID
 */
func (svc *UserService) GetUserById(id int64) (*User, error) {
	query := fmt.Sprintf("SELECT Username, TeamId, Status, Created from %s where Id = %d", USERS_TABLE, id)
	row := svc.DB.QueryRow(query)

	var user User
	var teamId int64
	err := row.Scan(&user.Username, &teamId, &user.Status, &user.Created)
	user.Id = id
	if err != nil {
		return nil, err
	}
	user.Id = id
	user.Team, err = svc.SG.TeamService.GetTeamById(teamId)
	return &user, err
}

/**
 * Get a user by username in a particular team.
 */
func (svc *UserService) GetUser(username string, team *Team) (*User, error) {
	query := fmt.Sprintf("SELECT Id, Status, Created from %s where Username = '%s' and TeamId = %d", USERS_TABLE, username, team.Id)
	row := svc.DB.QueryRow(query)

	var user User
	err := row.Scan(&user.Id, &user.Status, &user.Created)
	if err != nil {
		return nil, err
	}
	user.Username = username
	user.Team = team
	return &user, err
}

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
func (svc *UserService) SaveUser(user *User, override bool) error {
	if user.Id == 0 {
		id := UUIDGen()
		err := InsertRow(svc.DB, USERS_TABLE,
			"Id", id,
			"TeamId", user.Team.Id,
			"Username", user.Username,
			"Status", user.Status)
		if err == nil {
			user.Id = id
		} else {
			log.Println("Insert error.  Should retry: ", err)
		}
		return err
	} else {
		err := UpdateRows(svc.DB, USERS_TABLE, fmt.Sprintf("Id = %d", user.Id),
			"TeamId", user.Team.Id,
			"Username", user.Username,
			"Status", user.Status)
		if err.Error() == "No rows found" {
			// then Insert
			err = InsertRow(svc.DB, USERS_TABLE,
				"Id", user.Id,
				"TeamId", user.Team.Id,
				"Username", user.Username,
				"Status", user.Status)
		}
		return err
	}
}
