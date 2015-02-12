package services

import (
	. "github.com/panyam/backbone/models"
	. "gopkg.in/check.v1"
	// "code.google.com/p/gomock/gomock"
	// "log"
	// "time"
)

func CreateMessageService() IMessageService {
	return NewMessageService()
}

func (s *TestSuite) TestCreateMessageService(c *C) {
	svc := CreateMessageService()
	c.Assert(svc, Not(Equals), nil)
}

func (s *TestSuite) TestGetMessages(c *C) {
	chsvc := CreateChannelService()
	msgsvc := CreateMessageService()
	channel, _ := chsvc.CreateChannel("group", "test")
	msgs, _ := msgsvc.GetMessages(channel, nil, 0, -1)
	c.Assert(len(msgs), Equals, 0)
}

func (s *TestSuite) TestCreateMessage(c *C) {
	usersvc := CreateUserService()
	chsvc := CreateChannelService()
	msgsvc := CreateMessageService()
	sender, _ := usersvc.CreateUser("1", "user1")
	channel, _ := chsvc.CreateChannel("group", "test")
	message := NewMessage(channel, sender, nil)
	err := msgsvc.CreateMessage(message)
	c.Assert(err, Equals, nil)
	msgs, _ := msgsvc.GetMessages(channel, nil, 0, -1)
	c.Assert(len(msgs), Equals, 1)
}

func (s *TestSuite) TestDeleteMessage(c *C) {
	usersvc := CreateUserService()
	chsvc := CreateChannelService()
	msgsvc := CreateMessageService()
	sender, _ := usersvc.CreateUser("1", "user1")
	channel, _ := chsvc.CreateChannel("group", "test")
	message := NewMessage(channel, sender, nil)
	err := msgsvc.CreateMessage(message)
	c.Assert(err, Equals, nil)
	msgs, _ := msgsvc.GetMessages(channel, nil, 0, -1)
	c.Assert(len(msgs), Equals, 1)
	msgsvc.DeleteMessage(message)
	msgs, _ = msgsvc.GetMessages(channel, nil, 0, -1)
	c.Assert(len(msgs), Equals, 0)
}
