package models

import (
	"time"
)

/**
 * Channels/Threads/Groups etc
 */
type Channel struct {
	/**
	 * Globally unique Channel ID.
	 */
	Id string

	/**
	 * Name/Label of this channel.
	 */
	Name string

	/**
	 * Channel creation time.
	 */
	Created time.Time

	/**
	 * When the last message was posted on this channel.
	 */
	LastMessageAt time.Time

	/**
	 * Number of users in this channel.
	 */
	NumUsers int

	/**
	 * Status of this channel.
	 */
	Status string

	// A channel can be created by forking out of a message (like a seperate thread)
	Parent *Message

	// The team to which channel belongs.  Is this required?  or should it just
	// be created by an owner through which we can get the team.  The advantage
	// of linking to a user rather than a team is that if a user has multiple
	// teams he/she belongs to, then the channel can be allowed to be accessed
	// by members of those teams.

	// Alternatively the reason a team is required is to have a uniqueness
	// constraint on the channel name.  Is this constraint so important?  What
	// would happen if two channels have the same name?  Then the only issue
	// will be in identifying them by name - again not a big deal.  A team
	// may be too restrictive if we also had team/org break down - in which case
	// sometimes we may want to have a channel unique either within an org
	// or a team - so why not just a channel group and make that unique
	// so any entity that has this group will be the "tie" for channels it holds
	// By having a channel group any body can create a channel group (be it an
	// org or a team and have all channels under that group be unique without
	// worrying about the ownership structure and domain of the channel
	Group string

	// Metadata associated with the channel
	MetaData map[string]interface{}
}

type Message struct {
	/**
	 * Message ID - GUID
	 */
	Id string

	/**
	 * Channel to which this message belongs.
	 */
	Channel *Channel

	/**
	 * Sender of the message
	 */
	Sender *User

	/**
	 * When the message was sent.
	 */
	SentAt time.Time

	/**
	 * Type of message - eg, "invite", "status", "command", "event" etc
	 * This is normally only required when we do integrations so these
	 * messages could be sent by bots or other APIs
	 */
	MsgType string

	/**
	 * Whether message is to be persisted or not.  When a message is not
	 * persisted it is only sent to the users that are reachable.
	 */
	Persist bool

	/**
	 * All the message fragments.
	 */
	Fragments []*MessageFragment

	// Metadata associated with the message
	MetaData map[string]interface{}
}

/**
 * Keeps track of a message between the sender and the receiver.
 */
type MessageReceipt struct {
	Status     int
	Receiver   *User
	ReceivedAt time.Time
	Error      error
}

type MessageFragment struct {
	MimeType string
	Body     []byte
	Size     int64
}

func NewMessage(channel *Channel, sender *User) *Message {
	msg := Message{
		Channel: channel,
		Sender:  sender,
		SentAt:  time.Now()}
	return &msg
}

func NewChannel(id string, group string, name string) *Channel {
	return &Channel{Id: id, Group: group, Name: name}
}
