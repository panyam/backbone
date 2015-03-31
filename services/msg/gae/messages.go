package gae

import (
	"appengine"
	"appengine/datastore"
	. "github.com/panyam/relay/services/msg/core"
	"log"
)

type MessageService struct {
	Cls               interface{}
	context           appengine.Context
	messageCounter    int64
	messagesInChannel map[int64][]*Message
}

func NewMessageService(ctx appengine.Context) *MessageService {
	svc := MessageService{}
	svc.Cls = &svc
	svc.context = ctx
	svc.messagesInChannel = make(map[int64][]*Message)
	return &svc
}

/**
 * Get the messages in a channel for a particular user.
 * If channel does not exist then an empty list is returned.
 * If user is nil then all messages in a channel is returned.
 * Pagination is possible with offset and count.
 */
func (m *MessageService) GetMessages(channel *Channel, user *User, offset int,
	count int) ([]*Message, error) {
	return m.messagesInChannel[channel.Id], nil
}

/**
 * Send messages to particular recipients in this channel.
 * If Channel, message or user is nil an error is returned.
 */
func (m *MessageService) CreateMessage(message *Message) error {
	msgs, ok := m.messagesInChannel[message.Channel.Id]
	if !ok {
		msgs = make([]*Message, 0, 10)
		m.messagesInChannel[message.Channel.Id] = msgs
	}
	msgs = append(msgs, message)
	m.messagesInChannel[message.Channel.Id] = msgs
	return nil
}

/**
 * Remove a particular message.
 */
func (m *MessageService) DeleteMessage(message *Message) error {
	msgs, ok := m.messagesInChannel[message.Channel.Id]
	if !ok {
		return nil
	}
	for i, item := range msgs {
		if item == message {
			m.messagesInChannel[message.Channel.Id] = append(msgs[:i], msgs[i+
				1:]...)
			break
		}
	}
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
 * Removes all messages.
 */
func (svc *MessageService) RemoveAllMessages() {
	q := datastore.NewQuery("Message").KeysOnly()
	keys, err := q.GetAll(svc.context, nil)
	if err != nil {
		log.Println("RemoveAll Error: ", err)
	}
	datastore.DeleteMulti(svc.context, keys)
}
