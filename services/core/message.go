package core

import (
	"time"
)

type Message struct {
	/**
	 * Message ID - GUID
	 */
	Id int64

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
	index    uint32
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
