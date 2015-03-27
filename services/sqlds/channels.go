package sqlds

import (
	"database/sql"
	"errors"
	"fmt"
	. "github.com/panyam/backbone/services/core"
)

type ChannelService struct {
	Cls interface{}
	DB  *sql.DB
	SG  *ServiceGroup
}

const CHANNELS_TABLE = "channels"

func NewChannelService(db *sql.DB, sg *ServiceGroup) *ChannelService {
	svc := ChannelService{}
	svc.Cls = &svc
	svc.SG = sg
	svc.DB = db
	svc.InitDB()
	return &svc
}

func (svc *ChannelService) InitDB() {
	CreateTable(svc.DB, CHANNELS_TABLE,
		[]string{
			"Id bigint PRIMARY KEY",
			"TeamId bigint NOT NULL REFERENCES teams (Id)",
			"UserId bigint NOT NULL REFERENCES users (Id)",
			"Created TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp()",
			"Name varchar(128) DEFAULT ('')",
			"GroupName varchar(128) DEFAULT ('')",
			"Status INT DEFAULT (0)",
		})
}

/**
 * Lets a user create a channel.
 */
func (svc *ChannelService) SaveChannel(channel *Channel, override bool) error {
	if channel.Id == 0 {
		id := UUIDGen()
		query := fmt.Sprintf(`INSERT INTO %s ( Id, TeamId, UserId, Name, GroupName, Status) VALUES (%d, %d, %d, '%s', '%s', %d)`, CHANNELS_TABLE, id, channel.Team.Id, channel.Creator.Id, channel.Name, channel.GroupName, channel.Status)
		_, err := svc.DB.Exec(query)
		if err == nil {
			channel.Id = id
		}
		return err
	} else {
		query := fmt.Sprintf(`UPDATE %s SET GroupName = '%s', TeamId = %d, UserId= %d, Name = '%s', Status = %d where GroupName = '%s', Id = %d`, CHANNELS_TABLE, channel.GroupName, channel.Team.Id, channel.Creator.Id, channel.Name, channel.Status, channel.Id)
		_, err := svc.DB.Exec(query)
		return err
	}
}

/**
 * Retrieve a channel by Name.
 */
func (svc *ChannelService) GetChannelById(id int64) (*Channel, error) {
	query := fmt.Sprintf("SELECT TeamId, UserId, Status, Created, Name, GroupName from %s where Id = %d", CHANNELS_TABLE, id)
	row := svc.DB.QueryRow(query)

	var channel Channel
	var teamId int64
	var userId int64
	err := row.Scan(&teamId, &userId, &channel.Status, &channel.Created, &channel.Name, &channel.GroupName)
	if err != nil {
		return nil, err
	}
	channel.Id = id
	channel.Team, err = svc.SG.TeamService.GetTeamById(teamId)
	channel.Creator, err = svc.SG.UserService.GetUserById(userId)
	return &channel, err
}

/**
 * Delete a channel.
 */
func (c *ChannelService) DeleteChannel(channel *Channel) error {
	return errors.New("No such channel")
}

/**
 * Lets a user to join a channel (if allowed)
 */
func (c *ChannelService) JoinChannel(channel *Channel, user *User) error {
	return nil
}

/**
 * Lets a user leave a channel or be kicked out.
 */
func (c *ChannelService) LeaveChannel(channel *Channel, user *User) error {
	return nil
}

/**
 * Returns the channels the user belongs to.
 */
func (c *ChannelService) ListChannels(user *User, team *Team) ([]*Channel, error) {
	return nil, nil
}

/**
 * Removes all entries.
 */
func (svc *ChannelService) RemoveAllChannels() {
	ClearTable(svc.DB, CHANNELS_TABLE)
}
