package core

import (
	. "github.com/panyam/relay/utils"
	"time"
)

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
	channel := &Channel{Team: team, Name: name, Creator: creator}
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
