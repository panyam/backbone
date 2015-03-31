package core

import (
	"time"
)

type Message struct {
	Object

	/**
	 * Channel to which this message belongs.
	 */
	Channel *Channel

	/**
	 * Sender of the message
	 */
	Sender *User

	/**
	 * Type of message - eg, "invite", "status", "command", "event" etc
	 * This is normally only required when we do integrations so these
	 * messages could be sent by bots or other APIs
	 */
	MsgType string

	/**
	 * For storing simple messages.  The alternative is to use message fragments
	 * and store them as a list of fragments.
	 */
	MsgData string
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
		Sender:  sender}
	msg.Object = Object{Id: 0}
	return &msg
}
