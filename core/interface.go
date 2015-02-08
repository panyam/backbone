package core

import (
	"time"
)

type User interface {
	UserId() string
	Username() string
	SetUsername(value string)
	Addresses() []Address
	AddAddress(address *Address)
	GetMetaData(key string) interface{}
	SetMetaData(key string, value interface{})
}

type Address interface {
	Label() string
	SetLabel(label string)

	/**
	 * Where the user is from - eg FB, github, Twitter, email, phone etc
	 */
	Domain() string

	/**
	 * User's address within the domain.
	 */
	ID() string
}

type Team interface {
	TeamId() string
	GetName() string
	Created() time.Time
	Status() int
}

type Channel interface {
	Team() *Team
	Type() string
	Name() string
	Created() time.Time
	Status() string

	// A channel can be created by forking out of a message (like a seperate thread)
	Parent() Message
}

type Message interface {
	/**
	 * Channel to which this message belongs.
	 */
	Channel() Channel

	/**
	 * Type of message - eg, "invite", "status", "command", "event" etc
	 */
	Type() string

	/**
	 * Sender of the message
	 */
	Sender() User

	/**
	 * All the message fragments.
	 */
	Fragments() []MessageFragment

	/**
	 * Receipts of the messages indicating status of the recipients of the message.
	 */
	Receipts() []MessageReceipt

	/**
	 * When the message was sent.
	 */
	Sent() time.Time
}

/**
 * Keeps track of a message between the sender and the receiver.
 */
type MessageReceipt interface {
	Receiver() User
	Status() int
	Received() time.Time
	Error() error
}

type MessageFragment interface {
	MimeType() string
	Body() []byte
	Size() int64
}
