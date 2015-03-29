package sqlds

import (
	"database/sql"
	"fmt"
	. "github.com/panyam/backbone/services/core"
)

type ChannelService struct {
	Cls interface{}
	DB  *sql.DB
	SG  *ServiceGroup
}

const CHANNELS_TABLE = "channels"
const CHANNEL_MEMBERS_TABLE = "channel_members"
const CHANNEL_METADATA_TABLE = "channel_metadata"

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
			"TeamId bigint NOT NULL REFERENCES teams (Id) ON DELETE CASCADE",
			"UserId bigint DEFAULT(0) REFERENCES users (Id) ON DELETE SET DEFAULT",
			"Created TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp()",
			"Name TEXT DEFAULT ('')",
			"GroupName TEXT DEFAULT ('')",
			"LastMessageAt TIMESTAMP WITHOUT TIME ZONE",
			"Status INT DEFAULT (0)",
		})

	CreateTable(svc.DB, CHANNEL_MEMBERS_TABLE,
		[]string{
			"UserId bigint NOT NULL REFERENCES users (Id) ON DELETE CASCADE",
			"ChannelId bigint NOT NULL REFERENCES channels (Id) ON DELETE CASCADE",
			"JoinedAt TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp()",
			"LeftAt TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL",
			"Status INT DEFAULT (0)",
		},
		", CONSTRAINT unique_channel_membership UNIQUE (UserId, ChannelId)")

	CreateMetadataTable(svc.DB, CHANNEL_METADATA_TABLE, CHANNELS_TABLE, "Channel")
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
		query := fmt.Sprintf(`UPDATE %s SET GroupName = '%s', TeamId = %d, UserId = %d, Name = '%s', Status = %d where  Id = %d`, CHANNELS_TABLE, channel.GroupName, channel.Team.Id, channel.Creator.Id, channel.Name, channel.Status, channel.Id)
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

func (svc *ChannelService) GetChannelMembers(channel *Channel) []ChannelMember {
	query := fmt.Sprintf("SELECT UserId, JoinedAt, LeftAt, Status FROM %s where ChannelId = %d", CHANNEL_MEMBERS_TABLE, channel.Id)
	rows, err := svc.DB.Query(query)
	if err == nil {
		defer rows.Close()
	}

	members := make([]ChannelMember, 0, 0)
	for rows.Next() {
		var member ChannelMember
		var userId int64
		rows.Scan(&userId, &member.JoinedAt, &member.LeftAt, &member.Status)
		member.User, _ = svc.SG.UserService.GetUserById(userId)
		members = append(members, member)
	}
	return members
}

/**
 * Delete a channel.
 */
func (svc *ChannelService) DeleteChannel(channel *Channel) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE Id = %d", CHANNELS_TABLE, channel.Id)
	_, err := svc.DB.Exec(query)
	return err
}

/**
 * Tells if a user belongs to a channel.
 */
func (svc *ChannelService) ContainsUser(channel *Channel, user *User) bool {
	members := svc.GetChannelMembers(channel)
	for _, value := range members {
		if value.User.Id == user.Id {
			return true
		}
	}
	return false
}

/**
 * Lets a user to join a channel (if allowed)
 */
func (svc *ChannelService) JoinChannel(channel *Channel, user *User) error {
	query := fmt.Sprintf(`INSERT INTO %s ( UserId, ChannelId ) VALUES (%d, %d)`,
		CHANNEL_MEMBERS_TABLE, user.Id, channel.Id)
	_, err := svc.DB.Exec(query)
	return err
}

/**
 * Lets a user leave a channel or be kicked out.
 */
func (svc *ChannelService) LeaveChannel(channel *Channel, user *User) error {
	query := fmt.Sprintf(`UPDATE %s set UserId = %d, ChannelId = %d LeftAt = timestamp_now()`, CHANNEL_MEMBERS_TABLE, user.Id, channel.Id)
	_, err := svc.DB.Exec(query)
	return err
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
