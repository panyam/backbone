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

func NewMessageService(db *sql.DB, sg *ServiceGroup) *MessageService {
	svc := MessageService{}
	svc.Cls = &svc
	svc.SG = sg
	svc.DB = db
	svc.RemoveAllMessages()
	return &svc
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
