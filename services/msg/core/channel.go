package core

import (
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
	return nil, nil
}
