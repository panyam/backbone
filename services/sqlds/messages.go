package sqlds

import (
	"database/sql"
	. "github.com/panyam/backbone/services/core"
)

type MessageService struct {
	Cls interface{}
	DB  *sql.DB
	SG  *ServiceGroup
}

const MESSAGES_TABLE = "messages"
const MESSAGE_RECEIPTS_TABLE = "message_receipts"
const MESSAGE_FRAGMENTS_TABLE = "message_fragments"

func NewMessageService(db *sql.DB, sg *ServiceGroup) *MessageService {
	svc := MessageService{}
	svc.Cls = &svc
	svc.SG = sg
	svc.DB = db
	svc.InitDB()
	return &svc
}

func (svc *MessageService) InitDB() {
	CreateTable(svc.DB, MESSAGES_TABLE,
		[]string{
			"Id bigint PRIMARY KEY",
			"ChannelId bigint NOT NULL REFERENCES channels (Id) ON DELETE CASCADE",
			"SenderId bigint DEFAULT(0) REFERENCES users (Id) ON DELETE SET DEFAULT",
			"Created TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp()",
			"MsgType TEXT DEFAULT ('')",
			"MsgData TEXT DEFAULT ('')",
			"Status INT DEFAULT (0)",
		})
	CreateTable(svc.DB, MESSAGE_FRAGMENTS_TABLE,
		[]string{
			"MessageId bigint NOT NULL REFERENCES messages (Id) ON DELETE CASCADE",
			"Index INT NOT NULL",
			"Size INT DEFAULT(0)",
			"MimeType TEXT DEFAULT ('')",
			"Body BYTEA DEFAULT('')",
		},
		", CONSTRAINT unique_message_fragment UNIQUE (MessageId, Index)")
	CreateTable(svc.DB, MESSAGE_RECEIPTS_TABLE,
		[]string{
			"MessageId bigint NOT NULL REFERENCES messages (Id) ON DELETE CASCADE",
			"ReceiverId bigint NOT NULL REFERENCES users (Id) ON DELETE CASCADE",
			"ReceivedAt TIMESTAMP WITHOUT TIME ZONE DEFAULT statement_timestamp()",
			"Error TEXT DEFAULT('')",
			"Status INT DEFAULT(0)",
		},
		", CONSTRAINT unique_message_receipt UNIQUE (MessageId, ReceiverId)")
}

/**
 * Get the messages in a channel for a particular user.
 * If channel does not exist then an empty list is returned.
 * If user is nil then all messages in a channel is returned.
 * Pagination is possible with offset and count.
 */
func (m *MessageService) GetMessages(channel *Channel, user *User, offset int, count int) ([]*Message, error) {
	return nil, nil
}

/**
 * Send messages to particular recipients in this channel.
 * If Channel, message or user is nil an error is returned.
 */
func (m *MessageService) CreateMessage(message *Message) error {
	return nil
}

/**
 * Remove a particular message.
 */
func (m *MessageService) DeleteMessage(message *Message) error {
	return nil
}

/**
 * Saves a message.
 * If the message ID is missing (or empty) then a new message is created.
 * If message ID is present then the existing message is updated.
 */
func (m *MessageService) SaveMessage(message *Message) error {
	return nil
}

/**
 * Removes all entries.
 */
func (svc *MessageService) RemoveAllMessages() {
}
