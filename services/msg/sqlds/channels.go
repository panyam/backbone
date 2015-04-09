package sqlds

import (
	"database/sql"
	"fmt"
	. "github.com/panyam/relay/services/msg/core"
	. "github.com/panyam/relay/utils"
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
			"Public BOOL DEFAULT (true)",
			"LastMessageAt TIMESTAMP WITHOUT TIME ZONE",
			"Status INT DEFAULT (0)",
		})

	CreateTable(svc.DB, CHANNEL_MEMBERS_TABLE,
		[]string{
			"ChannelId bigint NOT NULL REFERENCES channels (Id) ON DELETE CASCADE",
			"Username TEXT NOT NULL",
			"JoinedAt TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp()",
			"LeftAt TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL",
			"Status INT DEFAULT (0)",
		},
		", CONSTRAINT unique_channel_membership UNIQUE (ChannelId, Username)")

	CreateMetadataTable(svc.DB, CHANNEL_METADATA_TABLE, CHANNELS_TABLE, "Channel")
}

/**
 * Lets a user create a channel.
 */
func (svc *ChannelService) SaveChannel(channel *Channel, override bool) error {
	if channel.Id == 0 {
		id := UUIDGen()
		err := InsertRow(svc.DB, CHANNELS_TABLE,
			"Id", id,
			"TeamId", channel.Team.Id,
			"UserId", channel.Creator.Id,
			"Name", channel.Name,
			"Public", channel.Public,
			"Status", channel.Status)
		if err == nil {
			channel.Id = id
		}
		return err
	} else {
		return UpdateRows(svc.DB, CHANNELS_TABLE, fmt.Sprintf("Id = %d", channel.Id),
			"TeamId", channel.Team.Id,
			"UserId", channel.Creator.Id,
			"Name", channel.Name,
			"Public", channel.Public,
			"Status", channel.Status)
	}
}

/**
 * Retrieve a channel by Name.
 */
func (svc *ChannelService) GetChannelById(id int64) (*Channel, error) {
	query := fmt.Sprintf("SELECT TeamId, UserId, Public, Status, Created, Name from %s where Id = %d", CHANNELS_TABLE, id)
	row := svc.DB.QueryRow(query)

	var channel Channel
	var teamId int64
	var userId int64
	err := row.Scan(&teamId, &userId, &channel.Public, &channel.Status, &channel.Created, &channel.Name)
	if err != nil {
		return nil, err
	}
	channel.Id = id
	channel.Team, err = svc.SG.TeamService.GetTeamById(teamId)
	channel.Creator, err = svc.SG.UserService.GetUserById(userId)
	return &channel, err
}

func (svc *ChannelService) GetChannelMembers(channel *Channel) []ChannelMember {
	query := fmt.Sprintf("SELECT Username, JoinedAt, LeftAt, Status FROM %s where ChannelId = %d", CHANNEL_MEMBERS_TABLE, channel.Id)
	rows, err := svc.DB.Query(query)
	if err == nil {
		defer rows.Close()
	}

	members := make([]ChannelMember, 0, 0)
	for rows.Next() {
		var member ChannelMember
		var username string
		rows.Scan(&username, &member.JoinedAt, &member.LeftAt, &member.Status)
		member.User, _ = svc.SG.UserService.GetUser(username, channel.Team)
		members = append(members, member)
	}
	return members
}

/**
 * Adds users to a channel.
 */
func (svc *ChannelService) AddChannelMembers(channel *Channel, usernames []string) error {
	for _, username := range usernames {
		err := InsertRow(svc.DB, CHANNEL_MEMBERS_TABLE,
			"ChannelId", channel.Id,
			"Username", username,
			"Status", 0,
			"LeftAt", nil)
		if err != nil {
			// then update
			err = UpdateRows(svc.DB, CHANNEL_MEMBERS_TABLE, fmt.Sprintf("ChannelId = %d AND Username = '%s'", channel.Id, username), "Status", 0, "LeftAt", nil)
		}
	}
	return nil
}

/**
 * Delete a channel.
 */
func (svc *ChannelService) DeleteChannel(channel *Channel) error {
	return DeleteById(svc.DB, CHANNELS_TABLE, channel.Id)
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
