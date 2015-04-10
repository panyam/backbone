package sqlds

import (
	"database/sql"
	"errors"
	"fmt"
	authcore "github.com/panyam/relay/services/auth/core"
	msgcore "github.com/panyam/relay/services/msg/core"
	. "github.com/panyam/relay/utils"
)

type AuthService struct {
	Cls         interface{}
	DB          *sql.DB
	UserService msgcore.IUserService
	TeamService msgcore.ITeamService
}

const REGISTRATIONS_TABLE = "registrations"
const LOGINS_TABLE = "logins"

func NewAuthService(db *sql.DB, userService msgcore.IUserService, teamService msgcore.ITeamService) *AuthService {
	svc := AuthService{}
	svc.Cls = &svc
	svc.DB = db
	svc.UserService = userService
	svc.TeamService = teamService
	svc.InitDB()
	return &svc
}

func (svc *AuthService) InitDB() {
	svc.DB.QueryRow("CREATE SEQUENCE registrationids MINVALUE 1")
	svc.DB.QueryRow("CREATE SEQUENCE loginids MINVALUE 1")
	CreateTable(svc.DB, REGISTRATIONS_TABLE,
		[]string{
			"Id bigint PRIMARY KEY",
			"TeamId bigint NOT NULL REFERENCES teams (Id) ON DELETE CASCADE",
			"Username TEXT NOT NULL",
			"Created TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp()",
			"ExpiresAt TIMESTAMP WITHOUT TIME ZONE",
			"AddressType TEXT DEFAULT('phone')",
			"Address TEXT NOT NULL",
			"Status INT DEFAULT(0)",
			"VerificationData TEXT DEFAULT('')",
		},
		", CONSTRAINT unique_registrations UNIQUE (TeamId, Username)")
	CreateTable(svc.DB, LOGINS_TABLE,
		[]string{
			"Id bigint PRIMARY KEY",
			"LoginType TEXT NOT NULL",
			"LoginToken TEXT NOT NULL",
			"UserId bigint NOT NULL REFERENCES users (Id) ON DELETE CASCADE",
			"Credentials TEXT DEFAULT('')",
			"Status INT DEFAULT(0)",
			"Created TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp()",
		},
		", CONSTRAINT unique_user_logins UNIQUE (LoginType, LoginToken)")
}

/**
 * Registers a user in a team.  If the username already exists in the team,
 * then an error is returned.  Otherwise a registration object is returned.
 */
func (svc *AuthService) SaveRegistration(registration *authcore.Registration) error {
	user, err := svc.UserService.GetUser(registration.Username, registration.Team)
	if user != nil || err == nil {
		return errors.New(fmt.Sprintf("Username (%s) already exists in team", registration.Username))
	}

	// otherwise create a registration object
	if registration.Id == 0 {
		id := UUIDGen()
		err := InsertRow(svc.DB, REGISTRATIONS_TABLE,
			"Id", id,
			"TeamId", registration.Team.Id,
			"Username", registration.Username,
			"Address", registration.Address,
			"AddressType", registration.AddressType,
			"Status", registration.Status)
		if err == nil {
			registration.Id = id
		}
		return err
	} else {
		return UpdateRows(svc.DB, REGISTRATIONS_TABLE, fmt.Sprintf("Id = %d", registration.Id),
			"TeamId", registration.Team.Id,
			"Username", registration.Username,
			"Address", registration.Address,
			"AddressType", registration.AddressType,
			"Status", registration.Status)
	}
}

/**
 * Gets registration by ID.
 */
func (svc *AuthService) GetRegistrationById(id int64) (*authcore.Registration, error) {
	query := fmt.Sprintf("SELECT Username, TeamId, Status, Created, AddressType, Address, VerificationData from %s where Id = %d", REGISTRATIONS_TABLE, id)
	row := svc.DB.QueryRow(query)

	var registration authcore.Registration
	var teamId int64
	err := row.Scan(&registration.Username, &teamId, &registration.Status, &registration.Created, &registration.AddressType, &registration.Address, &registration.VerificationData)
	if err != nil {
		return nil, err
	}
	registration.Id = id
	registration.Team, err = svc.TeamService.GetTeamById(teamId)
	return &registration, err
}

/**
 * Delete a particular registration.
 */
func (svc *AuthService) DeleteRegistration(registration *authcore.Registration) error {
	return DeleteById(svc.DB, REGISTRATIONS_TABLE, registration.Id)
}

/**
 * Removes all entries.
 */
func (svc *AuthService) RemoveAllRegistrations() {
	ClearTable(svc.DB, REGISTRATIONS_TABLE)
	svc.DB.QueryRow("DROP SEQUENCE registrationids")
}

/**
 * Removes all entries.
 */
func (svc *AuthService) RemoveAllLogins() {
	ClearTable(svc.DB, LOGINS_TABLE)
	svc.DB.QueryRow("DROP SEQUENCE loginids")
}

/**
 * Saves a user login object.
 */
func (svc *AuthService) SaveUserLogin(login *authcore.UserLogin) error {
	var userId int64 = 0
	if login.User != nil {
		userId = login.User.Id
	}
	if login.Id == 0 {
		id := UUIDGen()
		err := InsertRow(svc.DB, LOGINS_TABLE,
			"Id", id,
			"LoginType", login.LoginType,
			"LoginToken", login.LoginToken,
			"UserId", userId,
			"Credentials", login.Credentials,
			"Status", login.Status)
		if err == nil {
			login.Id = id
		}
		return err
	} else {
		return UpdateRows(svc.DB, LOGINS_TABLE, fmt.Sprintf("Id = %d", login.Id),
			"LoginType", login.LoginType,
			"LoginToken", login.LoginToken,
			"UserId", userId,
			"Credentials", login.Credentials,
			"Status", login.Status)
	}
}
