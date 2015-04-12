package sqlds

import (
	"database/sql"
	"fmt"
	. "github.com/panyam/relay/services/msg/core"
	. "github.com/panyam/relay/utils"
	"log"
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
			"UserId bigint NOT NULL REFERENCES users (Id) ON DELETE CASCADE",
			"JoinedAt TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp()",
			"LeftAt TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL",
			"Status INT DEFAULT (0)",
		},
		", CONSTRAINT unique_channel_membership UNIQUE (ChannelId, UserId)")

	CreateMetadataTable(svc.DB, CHANNEL_METADATA_TABLE, CHANNELS_TABLE, "Channel")
	svc.SG.IDService.CreateDomain("channelids", 1, 2)
}

/**
 * Lets a user create a channel.
 */
func (svc *ChannelService) SaveChannel(channel *Channel, override bool) error {
	if channel.Id == 0 {
		id, err := svc.SG.IDService.NextID("channelids")
		if err != nil {
			return err
		}
		err = InsertRow(svc.DB, CHANNELS_TABLE,
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
		err := UpdateRows(svc.DB, CHANNELS_TABLE, fmt.Sprintf("Id = %d", channel.Id),
			"TeamId", channel.Team.Id,
			"UserId", channel.Creator.Id,
			"Name", channel.Name,
			"Public", channel.Public,
			"Status", channel.Status)
		if err != nil && err.Error() == "No rows found" {
			err = InsertRow(svc.DB, CHANNELS_TABLE,
				"Id", channel.Id,
				"TeamId", channel.Team.Id,
				"UserId", channel.Creator.Id,
				"Name", channel.Name,
				"Public", channel.Public,
				"Status", channel.Status)
		}
		return err
	}
}

func (svc *ChannelService) GetChannels(team *Team, creator *User, orderBy string, participants []*User, matchAll bool) ([]*Channel, [][]ChannelMember) {
	query := "SELECT A.Id, A.Name, A.CreatorId, A.LastMessageAt, A.Public, A.Status, A.NumUsers FROM " +
		"( SELECT " +
		"C.Id as Id, C.Name as Name, C.UserId as CreatorId, " +
		"C.LastMessageAt as LastMessageAt, C.Public as Public, C.Status, " +
		"COUNT(M.UserId) as NumUsers " +
		"FROM channels C, channel_members M where C.Id = M.ChannelId "

	matchingParticipants := participants != nil && len(participants) > 0
	if matchingParticipants {
		query += "AND M.UserId in ("
		for index, part := range participants {
			if index > 0 {
				query += ", "
			}
			query += fmt.Sprintf("%s", FormatSqlValue(part.Id))
		}
		query += ")"
	}

	if creator != nil {
		query += fmt.Sprintf(" AND C.UserId = %d", creator.Id)
	}
	query += " AND M.Status = 0 GROUP BY (C.Id) ) A"

	if matchingParticipants && matchAll {
		query += fmt.Sprintf(" WHERE A.NumUsers = %d", len(participants))
	}

	if orderBy != "" {
		query += " ORDER BY C." + orderBy
	}

	log.Println("Query: ", query)
	rows, err := svc.DB.Query(query)
	if err == nil {
		defer rows.Close()
	}

	channels := make([]*Channel, 0, 0)
	members := make([][]ChannelMember, 0, 0)
	for rows.Next() {
		var creatorId int64
		var numUsers int = 10
		var lastMsgAt NullTime

		channel := Channel{Team: team, Creator: creator}
		channel.Team = team
		// id name creatorid lastmessageat public status numusers users statuses
		err := rows.Scan(&channel.Id, &channel.Name, &creatorId,
			&lastMsgAt, &channel.Public, &channel.Status, &numUsers)
		if err != nil {
			log.Println("Scan Error: ", err)
		}
		if creator == nil {
			channel.Creator, _ = svc.SG.UserService.GetUserById(creatorId)
		}
		channels = append(channels, &channel)
		members = append(members, svc.GetChannelMembers(&channel))
	}
	return channels, members
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
	query := fmt.Sprintf("SELECT U.Username, U.Id, M.JoinedAt, M.Status FROM %s M, %s U ", CHANNEL_MEMBERS_TABLE, USERS_TABLE)
	query += fmt.Sprintf(" WHERE ChannelId = %d AND M.Status = 0 AND M.UserId = U.Id", channel.Id)
	rows, err := svc.DB.Query(query)
	if err == nil {
		defer rows.Close()
	}

	members := make([]ChannelMember, 0, 0)
	for rows.Next() {
		var user = User{Team: channel.Team}
		var member = ChannelMember{User: &user}
		rows.Scan(&user.Username, &user.Id, &member.JoinedAt, &member.Status)
		members = append(members, member)
	}
	return members
}

/**
 * Adds users to a channel.
 */
func (svc *ChannelService) AddChannelMembers(channel *Channel, usernames []string) error {
	for _, username := range usernames {
		user, err := svc.SG.UserService.GetUser(username, channel.Team)
		if err != nil {
			// then create it
			user = NewUser(0, username, channel.Team)
			svc.SG.UserService.SaveUser(user, false)
		}
		err = InsertRow(svc.DB, CHANNEL_MEMBERS_TABLE,
			"ChannelId", channel.Id,
			"UserId", user.Id,
			"Status", 0,
			"LeftAt", nil)
		if err != nil {
			// then update
			err = UpdateRows(svc.DB, CHANNEL_MEMBERS_TABLE, fmt.Sprintf("ChannelId = %d AND UserId = %d", channel.Id, user.Id), "Status", 0, "LeftAt", nil)
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
	log.Println("Members in test: ", members)
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
	query := fmt.Sprintf(`INSERT INTO %s ( ChannelId, UserId) VALUES (%d, %d)`,
		CHANNEL_MEMBERS_TABLE, channel.Id, user.Id)
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
 * Removes all entries.
 */
func (svc *ChannelService) RemoveAllChannels() {
	ClearTable(svc.DB, CHANNELS_TABLE)
}
