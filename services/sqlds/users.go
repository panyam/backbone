package sqlds

import (
	"database/sql"
	"fmt"
	. "github.com/panyam/backbone/services/core"
)

type UserService struct {
	Cls interface{}
	SG  ServiceGroup
	DB  *sql.DB
}

const USERS_TABLE = "users"

func NewUserService(db *sql.DB) *UserService {
	svc := UserService{}
	svc.Cls = &svc
	svc.DB = db
	svc.InitDB()
	return &svc
}

func (svc *UserService) InitDB() {
	CreateTable(svc.DB, USERS_TABLE,
		[]string{
			"Id int64 PRIMARY KEY",
			"Username varchar(32) NOT NULL",
			"TeamId varchar(16) NOT NULL REFERENCES teams (Id)",
			"Created TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp()",
			"Status INT DEFAULT (0)",
		}, "UNIQUE (Username, TeamId)")
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
	query := fmt.Sprintf("SELECT Id, Username, TeamId, Status, Created from %s where Id = %d", USERS_TABLE, id)
	rows, err := svc.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rows.Next()

	var user User
	user.Id = id
	err = rows.Scan(&id)
	if err != nil {
		return nil, err
	}
	err = rows.Scan(&user.Username)
	if err != nil {
		return nil, err
	}
	var teamId int64
	err = rows.Scan(teamId)
	if err != nil {
		return nil, err
	}
	user.Team, err = svc.SG.TeamService.GetTeamById(teamId)

	err = rows.Scan(&user.Status)
	if err != nil {
		return nil, err
	}
	err = rows.Scan(&user.Created)
	if err != nil {
		return nil, err
	}
	return &user, err
}

/**
 * Get a user by username in a particular team.
 */
func (svc *UserService) GetUser(username string, team *Team) (*User, error) {
	return nil, nil
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
		query := fmt.Sprintf(`INSERT INTO %s ( TeamId, Username, Status) VALUES (%d, '%s')`, USERS_TABLE, user.Team.Id, user.Username, user.Status)
		_, err := svc.DB.Exec(query)
		return err
	} else {
		query := fmt.Sprintf(`UPDATE %s SET ( TeamId = %d, Username = '%s', Status = %d)`, USERS_TABLE, user.Team.Id, user.Username, user.Status)
		_, err := svc.DB.Exec(query)
		return err
	}
}
