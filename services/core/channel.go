package core

import (
	"time"
)

/**
 * Channels/Threads/Groups/Conversations etc
 */
type Channel struct {
	/**
	 * Globally unique Channel ID.
	 */
	Id int64

	/**
	 * Team this channel belongs to.   Mandatory.
	 */
	Team *Team

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
	Group string

	/**
	 * Creator of the group.
	 */
	Creator *User

	/**
	 * Channel creation time.
	 */
	Created time.Time

	/**
	 * When the last message was posted on this channel.
	 */
	LastMessageAt time.Time

	/**
	 * List of participants in this group.
	 */
	Participants []*User

	/**
	 * Status of this channel.
	 */
	Status string

	// A channel can be created by forking out of a message (like a seperate thread)
	ParentMessage *Message

	// Metadata associated with the channel
	MetaData map[string]interface{}
}

func NewChannel(team *Team, id int64, name string, group string) *Channel {
	return &Channel{Team: team, Id: id, Name: name, Group: group}
}

/**
 * Tells if a user belongs to a channel.
 */
func (channel *Channel) ContainsUser(user *User) bool {
	for _, value := range channel.Participants {
		if value.Id == user.Id {
			return true
		}
	}
	return false
}
