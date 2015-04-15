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
	svc.SG.IDService.CreateDomain(&CreateDomainRequest{nil, "userids", 1, 2})
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
func (svc *UserService) RemoveAllUsers(request *Request) {
	ClearTable(svc.DB, USERS_TABLE)
}

/**
 * Get user info by ID
 */
func (svc *UserService) GetUserById(user *User) (*User, error) {
	query := fmt.Sprintf("SELECT Username, TeamId, Status, Created from %s where Id = %d", USERS_TABLE, user.Id)
	row := svc.DB.QueryRow(query)

	var teamId int64
	err := row.Scan(&user.Username, &teamId, &user.Status, &user.Created)
	if err != nil {
		return nil, err
	}
	user.Team, err = svc.SG.TeamService.GetTeam(NewTeam(teamId, "", ""))
	if err == nil {
		user.Loaded = true
	}
	return user, err
}

/**
 * Get a user by username in a particular team.
 */
func (svc *UserService) GetUser(user *User) (*User, error) {
	if user.Id != 0 {
		return svc.GetUserById(user)
	}

	query := fmt.Sprintf("SELECT Id, Status, Created from %s where Username = '%s' and TeamId = %d", USERS_TABLE, user.Username, user.Team.Id)
	row := svc.DB.QueryRow(query)

	err := row.Scan(&user.Id, &user.Status, &user.Created)
	if err != nil {
		return nil, err
	}
	if err == nil {
		user.Loaded = true
	}
	return user, err
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
func (svc *UserService) SaveUser(request *SaveUserRequest) error {
	if request.User.Id == 0 {
		id, err := svc.SG.IDService.NextID(&NextIDRequest{nil, "userids"})
		if err != nil {
			return err
		}
		err = InsertRow(svc.DB, USERS_TABLE,
			"Id", id,
			"TeamId", request.User.Team.Id,
			"Username", request.User.Username,
			"Status", request.User.Status)
		if err == nil {
			request.User.Id = id
		} else {
			log.Println("Insert error.  Should retry: ", err)
		}
		return err
	} else {
		err := UpdateRows(svc.DB, USERS_TABLE, fmt.Sprintf("Id = %d", request.User.Id),
			"TeamId", request.User.Team.Id,
			"Username", request.User.Username,
			"Status", request.User.Status)
		if err.Error() == "No rows found" {
			// then Insert
			err = InsertRow(svc.DB, USERS_TABLE,
				"Id", request.User.Id,
				"TeamId", request.User.Team.Id,
				"Username", request.User.Username,
				"Status", request.User.Status)
		}
		return err
	}
}
