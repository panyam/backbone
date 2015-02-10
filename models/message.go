package models

import (
	"time"
)

/**
 * Channels/Threads/Groups etc
 */
type Channel struct {
	Team        *Team
	Id          string
	ChannelType string
	Name        string
	Created     time.Time
	Status      string

	// A channel can be created by forking out of a message (like a seperate thread)
	Parent *Message
}

type Message struct {
	/**
	 * Channel to which this message belongs.
	 */
	Channel *Channel

	/**
	 * Whether the message is in reply to another message.
	 */
	InReplyTo *Message

	/**
	 * Type of message - eg, "invite", "status", "command", "event" etc
	 */
	MsgType string

	/**
	 * Sender of the message
	 */
	Sender *User

	/**
	 * All the message fragments.
	 */
	Fragments []*MessageFragment

	/**
	 * Receipts of the messages indicating status of the recipients of the message.
	 */
	Receipts []*MessageReceipt

	/**
	 * When the message was sent.
	 */
	SentAt time.Time
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

func NewMessage() (*Message, error) {
	msg := Message{}
	return &msg, nil
}
