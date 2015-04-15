package core

import (
	. "github.com/panyam/relay/utils"
	"time"
)

type CreateChannelRequest struct {
	Channel  *Channel
	Override bool
}

type GetChannelsRequest struct {
	Team         *Team
	Creator      *User
	OrderBy      string
	Participants []*User
	MatchAll     bool
}

type GetChannelsResult struct {
	Channels []*Channel
	Members  [][]*ChannelMember
}

type GetChannelResult struct {
	Channel *Channel
	Members []*ChannelMember
}

type InviteMembersRequest struct {
	*Request
	Channel   *Channel
	Usernames []string
}

type ChannelMembershipRequest struct {
	*Request
	Channel *Channel
	User    *User
}

type DeleteChannelRequest struct {
	*Request
	Channel *Channel
}

type IChannelService interface {
	/**
	 * Removes all entries.
	 */
	RemoveAllChannels(request *Request)

	/**
	 * Creates a new channel.
	 */
	CreateChannel(request *CreateChannelRequest) (*Channel, error)

	/**
	 * Updates an existing channel.
	 */
	UpdateChannel(channel *Channel) (*Channel, error)

	/**
	 * Get channels meeting particular criterea
	 */
	GetChannels(request *GetChannelsRequest) (*GetChannelsResult, error)

	/**
	 * Get channel by Id
	 */
	GetChannelById(channelId int64) (*GetChannelResult, error)

	/**
	 * Adds users to a channel.
	 */
	AddChannelMembers(request *InviteMembersRequest) error

	/**
	 * Delete a channel.
	 */
	DeleteChannel(request *DeleteChannelRequest) error

	/**
	 * Tells if a channe contains a user.
	 */
	ContainsUser(request *ChannelMembershipRequest) bool

	/**
	 * Lets a user to join a channel (if allowed)
	 */
	JoinChannel(request *ChannelMembershipRequest) error

	/**
	 * Lets a user leave a channel or be kicked out.
	 */
	LeaveChannel(request *ChannelMembershipRequest) error
}

/**
 * Channels/Threads/Groups/Conversations etc
 */
type Channel struct {
	Object

	/**
	 * Team this channel belongs to.   Mandatory.
	 */
	Team *Team

	/**
	 * Creator of the group.
	 */
	Creator *User

	/**
	 * Name of this channel.
	 */
	Name string

	/**
	 * Is this a public a channel or only visible to members?
	 */
	Public bool

	/**
	 * When the last message was posted on this channel.
	 */
	LastMessageAt time.Time
}

type ChannelMember struct {
	User *User

	JoinedAt time.Time

	LeftAt time.Time

	Status int
}

func NewChannel(team *Team, creator *User, id int64, name string) *Channel {
	channel := &Channel{Team: team, Name: name, Creator: creator, Public: true}
	channel.Object = Object{Id: id}
	return channel
}

func ChannelFromDict(data map[string]interface{}) (*Channel, error) {
	channel := Channel{}
	channel.Name = data["Name"].(string)
	channel.Public = data["Public"].(bool)
	channel.Status, _ = JsonNumberToInt(data["Status"])
	channel.Id = String2ID(data["Id"].(string))
	return &channel, nil
}

func (c *Channel) ToDict() map[string]interface{} {
	out := map[string]interface{}{
		"Id":     ID2String(c.Id),
		"Name":   c.Name,
		"Public": c.Public,
		"Status": c.Status,
		"TeamId": ID2String(c.Team.Id),
		"UserId": ID2String(c.Creator.Id),
	}
	return out
}
