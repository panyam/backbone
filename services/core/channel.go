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
	 * Uniqueness groups allows a channel to be unique within a group.  No two
	 * channels can have the same uniqueness group. The point of this is so that
	 * the user can define uniqueness constraints on when a group should not be
	 * duplicated.
	 */
	GroupName string

	/**
	 * When the last message was posted on this channel.
	 */
	LastMessageAt time.Time

	// A channel can be created by forking out of a message (like a seperate thread)
	ParentMessage *Message

	// Members in this channel
	Members *[]ChannelMember

	// Metadata associated with the channel
	MetaData map[string]interface{}
}

type ChannelMember struct {
	User *User

	JoinedAt time.Time

	LeftAt time.Time

	Status int
}

func NewChannel(team *Team, creator *User, id int64, name string, group string) *Channel {
	channel := &Channel{Team: team, Name: name, GroupName: group, Creator: creator}
	channel.Object = Object{Id: id}
	return channel
}

/**
 * Tells if a user belongs to a channel.
 */
func (channel *Channel) ContainsUser(user *User) bool {
	if channel.Members != nil {
		for _, value := range *channel.Members {
			if value.User.Id == user.Id {
				return true
			}
		}
	}
	return false
}
