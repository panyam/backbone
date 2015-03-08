package services

import (
	. "github.com/panyam/backbone/models"
	. "gopkg.in/check.v1"
	// "code.google.com/p/gomock/gomock"
	// "log"
	// "time"
)

func (s *TestSuite) TestCreateMessageService(c *C) {
	svc := CreateMessageService()
	c.Assert(svc, Not(Equals), nil)
}

func (s *TestSuite) TestGetMessages(c *C) {
	chsvc := CreateChannelService()
	team := NewTeam("superteam", "superorg", "Super Team")
	channel := NewChannel(team, "", "test", "group")
	chsvc.SaveChannel(channel, true)

	msgsvc := CreateMessageService()
	msgs, _ := msgsvc.GetMessages(channel, nil, 0, -1)
	c.Assert(len(msgs), Equals, 0)
}

func (s *TestSuite) TestCreateMessage(c *C) {
	chsvc := CreateChannelService()
	team := NewTeam("superteam", "superorg", "Super Team")
	channel := NewChannel(team, "", "test", "group")
	err := chsvc.SaveChannel(channel, true)

	usersvc := CreateUserService()
	msgsvc := CreateMessageService()
	sender, _ := usersvc.CreateUser("1", "user1")
	message := NewMessage(channel, sender)
	err = msgsvc.CreateMessage(message)
	c.Assert(err, Equals, nil)
	msgs, _ := msgsvc.GetMessages(channel, nil, 0, -1)
	c.Assert(len(msgs), Equals, 1)
}

func (s *TestSuite) TestDeleteMessage(c *C) {
	chsvc := CreateChannelService()
	team := NewTeam("superteam", "superorg", "Super Team")
	channel := NewChannel(team, "", "test", "group")
	err := chsvc.SaveChannel(channel, true)

	usersvc := CreateUserService()
	msgsvc := CreateMessageService()
	sender, _ := usersvc.CreateUser("1", "user1")
	message := NewMessage(channel, sender)
	err = msgsvc.CreateMessage(message)
	c.Assert(err, Equals, nil)
	msgs, _ := msgsvc.GetMessages(channel, nil, 0, -1)
	c.Assert(len(msgs), Equals, 1)
	msgsvc.DeleteMessage(message)
	msgs, _ = msgsvc.GetMessages(channel, nil, 0, -1)
	c.Assert(len(msgs), Equals, 0)
}
