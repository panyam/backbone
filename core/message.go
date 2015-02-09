package core

import (
	"time"
)

/**
 * Channels/Threads/Groups etc
 */
type Channel struct {
	Cls         IChannel
	team        ITeam
	channelType string
	name        string
	created     time.Time
	status      string

	// A channel can be created by forking out of a message (like a seperate thread)
	parent *Message
}

type Message struct {
	Cls IMessage

	/**
	 * Channel to which this message belongs.
	 */
	channel *Channel

	/**
	 * Whether the message is in reply to another message.
	 */
	inReplyTo *Message

	/**
	 * Type of message - eg, "invite", "status", "command", "event" etc
	 */
	msgType string

	/**
	 * Sender of the message
	 */
	sender *User

	/**
	 * All the message fragments.
	 */
	fragments []*MessageFragment

	/**
	 * Receipts of the messages indicating status of the recipients of the message.
	 */
	receipts []*MessageReceipt

	/**
	 * When the message was sent.
	 */
	sentAt time.Time
}

/**
 * Keeps track of a message between the sender and the receiver.
 */
type MessageReceipt struct {
	status     int
	receiver   *User
	receivedAt time.Time
	error      error
}

type MessageFragment struct {
	mimeType string
	body     []byte
	size     int64
}

func (c *Message) Self() IMessage {
	return c.Cls
}

func NewMessage() (*Message, error) {
	msg := Message{}
	msg.Cls = &msg
	return &msg, nil
}

func (m *Message) Channel() IChannel {
	return m.channel
}

func (m *Message) Type() string {
	return m.msgType
}

func (m *Message) Sender() IUser {
	return m.sender
}

func (m *Message) NumFragments() int {
	return len(m.fragments)
}

func (m *Message) Fragment(index int) IMessageFragment {
	return m.fragments[index]
}

func (m *Message) NumReceipts() int {
	return len(m.receipts)
}

func (m *Message) Receipt(index int) IMessageReceipt {
	return m.receipts[index]
}

func (m *Message) Sent() time.Time {
	return m.sentAt
}

func (m *MessageReceipt) Receiver() IUser {
	return m.receiver
}

func (m *MessageReceipt) Status() int {
	return m.status
}

func (m *MessageReceipt) Received() time.Time {
	return m.receivedAt
}

func (m *MessageReceipt) Error() error {
	return m.error
}

func (m *MessageFragment) MimeType() string {
	return m.mimeType
}

func (m *MessageFragment) Body() []byte {
	return m.body
}

func (m *MessageFragment) Size() int64 {
	return m.size
}

////////// Channel related methods

func (c *Channel) Team() ITeam {
	return c.team
}

func (c *Channel) Type() string {
	return c.channelType
}

func (c *Channel) Name() string {
	return c.name
}

func (c *Channel) Created() time.Time {
	return c.created
}

func (c *Channel) Status() string {
	return c.status
}

func (c *Channel) Parent() IMessage {
	return c.parent
}
