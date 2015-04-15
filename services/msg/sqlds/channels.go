package sqlds

import (
	"database/sql"
	"errors"
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
	request := &CreateDomainRequest{nil, "channelids", 1, 2}
	svc.SG.IDService.CreateDomain(request)
}

/**
 * Lets a user create a channel.
 */
func (svc *ChannelService) CreateChannel(request *CreateChannelRequest) (*Channel, error) {
	nextid_request := &NextIDRequest{nil, "channelids"}
	id, err := svc.SG.IDService.NextID(nextid_request)
	if err != nil {
		return nil, err
	}
	err = InsertRow(svc.DB, CHANNELS_TABLE,
		"Id", id,
		"TeamId", request.Channel.Team.Id,
		"UserId", request.Channel.Creator.Id,
		"Name", request.Channel.Name,
		"Public", request.Channel.Public,
		"Status", request.Channel.Status)
	if err == nil {
		request.Channel.Id = id
	}
	return request.Channel, err
}

func (svc *ChannelService) UpdateChannel(channel *Channel) (*Channel, error) {
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
	return channel, err
}

func (svc *ChannelService) GetChannels(request *GetChannelsRequest) (*GetChannelsResult, error) {
	if request.Team == nil {
		return nil, errors.New("Cannot filter channels without team")
	}

	if !request.Team.Loaded {
		request.Team, _ = svc.SG.TeamService.GetTeam(request.Team)
		if request.Team == nil {
			return nil, errors.New("No such team")
		}
	}

	// Verify owner/creator if any
	if request.Creator != nil {
		if !request.Creator.Loaded {
			request.Creator, _ = svc.SG.UserService.GetUser(request.Creator)
			if request.Creator == nil {
				return nil, errors.New("No such user")
			}
		}
	}

	// check participants
	newParticipants := make([]*User, 0, 0)
	for _, participant := range request.Participants {
		if !participant.Loaded {
			// try to load it
			participant, _ = svc.SG.UserService.GetUser(participant)
			if participant == nil {
				if request.MatchAll {
					// invalid participant means we cant find any channels
					return nil, nil
				}
			} else {
				newParticipants = append(newParticipants, participant)
			}
		}
	}
	request.Participants = newParticipants

	/*
		teamIdParam := mux.Vars(request)["teamId"]
		ownerParam := request.FormValue("owner")
		matchall := request.FormValue("matchall") == "true"
		participantsParam := strings.Split(request.FormValue("participants"), ",")
	*/

	// Begin the queries!
	query := "SELECT A.Id, A.Name, A.CreatorId, A.LastMessageAt, A.Public, A.Status, A.NumUsers FROM " +
		"( SELECT " +
		"C.Id as Id, C.Name as Name, C.UserId as CreatorId, " +
		"C.LastMessageAt as LastMessageAt, C.Public as Public, C.Status, " +
		"COUNT(M.UserId) as NumUsers " +
		"FROM channels C, channel_members M where C.Id = M.ChannelId "

	matchingParticipants := request.Participants != nil && len(request.Participants) > 0
	if matchingParticipants {
		query += "AND M.UserId in ("
		for index, part := range request.Participants {
			if index > 0 {
				query += ", "
			}
			query += fmt.Sprintf("%s", FormatSqlValue(part.Id))
		}
		query += ")"
	}

	if request.Creator != nil {
		query += fmt.Sprintf(" AND C.UserId = %d", request.Creator.Id)
	}
	query += " AND M.Status = 0 GROUP BY (C.Id) ) A"

	if matchingParticipants && request.MatchAll {
		query += fmt.Sprintf(" WHERE A.NumUsers = %d", len(request.Participants))
	}

	if request.OrderBy != "" {
		query += " ORDER BY C." + request.OrderBy
	}

	log.Println("Query: ", query)
	rows, err := svc.DB.Query(query)
	if err == nil {
		defer rows.Close()
	}

	channels := make([]*Channel, 0, 0)
	members := make([][]*ChannelMember, 0, 0)
	for rows.Next() {
		var creatorId int64
		var numUsers int = 10
		var lastMsgAt NullTime

		channel := Channel{Team: request.Team, Creator: request.Creator}
		channel.Team = request.Team
		// id name creatorid lastmessageat public status numusers users statuses
		err := rows.Scan(&channel.Id, &channel.Name, &creatorId,
			&lastMsgAt, &channel.Public, &channel.Status, &numUsers)
		if err != nil {
			log.Println("Scan Error: ", err)
		}
		if request.Creator == nil {
			channel.Creator, _ = svc.SG.UserService.GetUser(NewUser(creatorId, "", nil))
		}
		channels = append(channels, &channel)
		members = append(members, svc.GetChannelMembers(&channel))
	}
	return &GetChannelsResult{channels, members}, nil
}

/**
 * Retrieve a channel by Name.
 */
func (svc *ChannelService) GetChannelById(channelId int64) (*GetChannelResult, error) {
	query := fmt.Sprintf("SELECT TeamId, UserId, Public, Status, Created, Name from %s where Id = %d", CHANNELS_TABLE, channelId)
	row := svc.DB.QueryRow(query)

	var channel Channel
	var teamId int64
	var userId int64
	err := row.Scan(&teamId, &userId, &channel.Public, &channel.Status, &channel.Created, &channel.Name)
	if err != nil {
		return nil, err
	}
	channel.Id = channelId
	channel.Team, err = svc.SG.TeamService.GetTeam(NewTeam(teamId, "", ""))
	channel.Creator, err = svc.SG.UserService.GetUser(NewUser(userId, "", nil))
	result := GetChannelResult{Channel: &channel,
		Members: svc.GetChannelMembers(&channel)}
	if err == nil {
		channel.Loaded = true
	}
	return &result, err
}

func (svc *ChannelService) GetChannelMembers(channel *Channel) []*ChannelMember {
	query := fmt.Sprintf("SELECT U.Username, U.Id, M.JoinedAt, M.Status FROM %s M, %s U ", CHANNEL_MEMBERS_TABLE, USERS_TABLE)
	query += fmt.Sprintf(" WHERE ChannelId = %d AND M.Status = 0 AND M.UserId = U.Id", channel.Id)
	rows, err := svc.DB.Query(query)
	if err == nil {
		defer rows.Close()
	}

	members := make([]*ChannelMember, 0, 0)
	for rows.Next() {
		var user = User{Team: channel.Team}
		var member = ChannelMember{User: &user}
		rows.Scan(&user.Username, &user.Id, &member.JoinedAt, &member.Status)
		members = append(members, &member)
	}
	return members
}

/**
 * Adds users to a channel.
 */
func (svc *ChannelService) AddChannelMembers(request *InviteMembersRequest) error {
	for _, username := range request.Usernames {
		user, err := svc.SG.UserService.GetUser(NewUser(0, username, request.Channel.Team))
		if err != nil {
			// then create it
			user = NewUser(0, username, request.Channel.Team)
			save_request := &SaveUserRequest{nil, user, false}
			svc.SG.UserService.SaveUser(save_request)
		}
		err = InsertRow(svc.DB, CHANNEL_MEMBERS_TABLE,
			"ChannelId", request.Channel.Id,
			"UserId", user.Id,
			"Status", 0,
			"LeftAt", nil)
		if err != nil {
			// then update
			err = UpdateRows(svc.DB, CHANNEL_MEMBERS_TABLE, fmt.Sprintf("ChannelId = %d AND UserId = %d", request.Channel.Id, user.Id), "Status", 0, "LeftAt", nil)
		}
	}
	return nil
}

/**
 * Delete a channel.
 */
func (svc *ChannelService) DeleteChannel(request *DeleteChannelRequest) error {
	return DeleteById(svc.DB, CHANNELS_TABLE, request.Channel.Id)
}

/**
 * Tells if a user belongs to a channel.
 */
func (svc *ChannelService) ContainsUser(request *ChannelMembershipRequest) bool {
	members := svc.GetChannelMembers(request.Channel)
	log.Println("Members in test: ", members)
	for _, value := range members {
		if value.User.Id == request.User.Id {
			return true
		}
	}
	return false
}

/**
 * Lets a user to join a channel (if allowed)
 */
func (svc *ChannelService) JoinChannel(request *ChannelMembershipRequest) error {
	query := fmt.Sprintf(`INSERT INTO %s ( ChannelId, UserId) VALUES (%d, %d)`,
		CHANNEL_MEMBERS_TABLE, request.Channel.Id, request.User.Id)
	_, err := svc.DB.Exec(query)
	return err
}

/**
 * Lets a user leave a channel or be kicked out.
 */
func (svc *ChannelService) LeaveChannel(request *ChannelMembershipRequest) error {
	query := fmt.Sprintf(`UPDATE %s set UserId = %d, ChannelId = %d LeftAt = timestamp_now()`, CHANNEL_MEMBERS_TABLE, request.User.Id, request.Channel.Id)
	_, err := svc.DB.Exec(query)
	return err
}

/**
 * Removes all entries.
 */
func (svc *ChannelService) RemoveAllChannels(request *Request) {
	ClearTable(svc.DB, CHANNELS_TABLE)
}
