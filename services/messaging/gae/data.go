package gae

import (
	"appengine/datastore"
	. "github.com/panyam/relay/services/messaging/core"
	"time"
)

/**
 * GAEChannel representation to get around unsupported types
 */
type GAEChannel struct {
	/**
	 * Globally unique Channel ID.
	 */
	Id *datastore.Key

	/**
	 * Team this channel belongs to.   Mandatory.
	 */
	TeamKey *datastore.Key

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
	 * Creator of the group.
	 */
	CreatorKey *datastore.Key

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
	ParticipantKeys []*datastore.Key

	/**
	 * Status of this channel.
	 */
	Status string

	// A channel can be created by forking out of a message (like a seperate thread)
	ParentMessageKey *datastore.Key
}

func (c *ChannelService) ToChannel(gc *GAEChannel, channel *Channel) {
	/*
		channel.Id = gc.Id.StringID()
		channel.Creator = gc.CreatorKey = UserKeyFor(c.context, channel.Creator.Id)
		gc.Name = channel.Name
		gc.GroupName = channel.GroupName
		gc.Created = channel.Created
		gc.LastMessageAt = channel.LastMessageAt
		gc.ParticipantKeys = nil
		for _, user := range channel.Participants {
			gc.ParticipantKeys = append(gc.ParticipantKeys, UserKeyFor(c.context, user.Id))
		}
		gc.Status = channel.Status
		gc.ParentMessageKey = MessageKeyFor(c.context, channel.ParentMessage.Id)
	*/
}

func (c *ChannelService) FromChannel(gc *GAEChannel, channel *Channel) {
	gc.Id = ChannelKeyFor(c.context, channel.Id)
	gc.TeamKey = TeamKeyFor(c.context, channel.Team.Id)
	gc.CreatorKey = UserKeyFor(c.context, channel.Creator.Id)
	gc.Name = channel.Name
	gc.GroupName = channel.GroupName
	gc.Created = channel.Created
	gc.LastMessageAt = channel.LastMessageAt
	gc.ParticipantKeys = nil
	for _, user := range channel.Participants {
		gc.ParticipantKeys = append(gc.ParticipantKeys, UserKeyFor(c.context, user.Id))
	}
	// gc.Status = channel.Status
	gc.ParentMessageKey = MessageKeyFor(c.context, channel.ParentMessage.Id)
}
