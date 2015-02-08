package core

import (
	"time"
)

/**
 * Channels/Threads/Groups etc
 */
type SimpleChannel struct {
	Cls         Channel
	team        Team
	channelType string
	name        string
	created     time.Time
	status      string

	// A channel can be created by forking out of a message (like a seperate thread)
	parent Message
}

type SimpleMessage struct {
	Cls Message

	/**
	 * Channel to which this message belongs.
	 */
	channel Channel

	/**
	 * Whether the message is in reply to another message.
	 */
	inReplyTo Message

	/**
	 * Type of message - eg, "invite", "status", "command", "event" etc
	 */
	msgType string

	/**
	 * Sender of the message
	 */
	sender User

	/**
	 * All the message fragments.
	 */
	fragments []MessageFragment

	/**
	 * Receipts of the messages indicating status of the recipients of the message.
	 */
	receipts []MessageReceipt

	/**
	 * When the message was sent.
	 */
	sentAt time.Time
}

/**
 * Keeps track of a message between the sender and the receiver.
 */
type SimpleMessageReceipt struct {
	Receiver *User
	Status   int
	Received time.Time
	Error    error
}

type SimpleMessageFragment struct {
	MimeType string
	Body     []byte
	Size     int64
}

func (c *SimpleMessage) Self() Message {
	return c.Cls
}

func NewMessage() (*SimpleMessage, error) {
	msg := SimpleMessage{}
	msg.Cls = &msg
	return &msg, nil
}

func (m *SimpleMessage) Channel() Channel {
	return m.channel
}

func (m *SimpleMessage) Type() string {
	return m.msgType
}

func (m *SimpleMessage) Sender() User {
	return m.sender
}

func (m *SimpleMessage) Fragments() []MessageFragment {
	return m.fragments
}

func (m *SimpleMessage) Receipts() []MessageReceipt {
	return m.receipts
}

func (m *SimpleMessage) Sent() time.Time {
	return m.sentAt
}
