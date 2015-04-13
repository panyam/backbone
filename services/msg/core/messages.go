package core

import (
	"time"
)

type GetMessagesRequest struct {
	*Request
	Channel *Channel
	Sender  *User
	Offset  int
	Count   int
}

type SaveMessageRequest struct {
	*Request

	Message *Message
}

type GetMessageRequest struct {
	*Request
	Id int64
}

type GetMessageResult struct {
	Message   *Message
	Receipts  []*MessageReceipt
	Fragments []*MessageFragment
}

type DeleteMessageRequest struct {
	*Request
	Message *Message
}

type IMessageService interface {
	/**
	 * Removes all entries.
	 */
	RemoveAllMessages(request *Request)

	/**
	 * Get the messages in a channel for a particular user.
	 */
	GetMessages(request *GetMessagesRequest) ([]*Message, error)

	/**
	 * Creates a message to particular recipients in this channel.  This is
	 * called "Create" instead of "Send" so as to not confuse with the delivery
	 * details.
	 * If message ID is empty then the backend can auto generate one if it is
	 * capable of doing so.
	 * A valid Message object on return WILL have a non empty ID if the backend can
	 * auto generate IDs
	 */
	SaveMessage(request *SaveMessageRequest) error

	/**
	 * Gets a message by ID
	 */
	GetMessageById(request *GetMessageRequest) (*GetMessageResult, error)

	/**
	 * Gets the fragments of a message.
	 */
	GetMessageFragments(request *GetMessageRequest) (*GetMessageResult, error)

	/**
	 * Get receipts of a particular message.
	 */
	GetMessageReceipts(request *GetMessageRequest) (*GetMessageResult, error)

	/**
	 * Remove a particular message.
	 */
	DeleteMessage(request *DeleteMessageRequest) error
}

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
