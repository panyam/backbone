package core

import (
	"time"
)

type IUser interface {
	UserId() string
	Username() string
	SetUsername(value string)
	Addresses() []IAddress
	AddAddress(address IAddress)
	GetMetaData(key string) interface{}
	SetMetaData(key string, value interface{})
}

type IAddress interface {
	Label() string
	SetLabel(label string)

	/**
	 * Where the user is from - eg FB, github, Twitter, email, phone etc
	 */
	Domain() string

	/**
	 * IUser's address within the domain.
	 */
	ID() string
}

type ITeam interface {
	TeamId() string
	GetName() string
	Created() time.Time
	Status() int
}

type IChannel interface {
	Team() ITeam
	Type() string
	Name() string
	Created() time.Time
	Status() string

	// A channel can be created by forking out of a message (like a seperate thread)
	Parent() IMessage
}

type IMessage interface {
	/**
	 * IChannel to which this message belongs.
	 */
	Channel() IChannel

	/**
	 * Type of message - eg, "invite", "status", "command", "event" etc
	 */
	Type() string

	/**
	 * Sender of the message
	 */
	Sender() IUser

	/**
	 * Fragment count and fragment at index.
	 */
	NumFragments() int
	Fragment(index int) IMessageFragment

	/**
	 * Receipt count and receipt at index.
	 */
	NumReceipts() int
	Receipt(index int) IMessageReceipt

	/**
	 * When the message was sent.
	 */
	Sent() time.Time
}

/**
 * Keeps track of a message between the sender and the receiver.
 */
type IMessageReceipt interface {
	Receiver() IUser
	Status() int
	Received() time.Time
	Error() error
}

type IMessageFragment interface {
	MimeType() string
	Body() []byte
	Size() int64
}
