package services

/*
import (
	. "github.com/panyam/backbone/services/core"
	. "gopkg.in/check.v1"
	// "code.google.com/p/gomock/gomock"
	// "log"
	// "time"
)

func (s *TestSuite) TestCreateMessageService(c *C) {
	svc := s.serviceGroup.MessageService
	c.Assert(svc, Not(Equals), nil)
}

func (s *TestSuite) TestGetMessages(c *C) {
	chsvc := s.serviceGroup.ChannelService
	team := NewTeam("superteam", "superorg", "Super Team")
	channel := NewChannel(team, "", "test", "group")
	chsvc.SaveChannel(channel, true)

	msgsvc := s.serviceGroup.MessageService
	msgs, _ := msgsvc.GetMessages(channel, nil, 0, -1)
	c.Assert(len(msgs), Equals, 0)
}

func (s *TestSuite) TestCreateMessage(c *C) {
	chsvc := s.serviceGroup.ChannelService
	team, _ := s.serviceGroup.TeamService.CreateTeam("1", "org", "team")
	channel := NewChannel(team, "", "test", "group")
	err := chsvc.SaveChannel(channel, true)

	usersvc := s.serviceGroup.UserService
	msgsvc := s.serviceGroup.MessageService

	sender := User{Username: "user1", Team: team}
	err := usersvc.SaveUser(&sender, false)

	message := NewMessage(channel, sender)
	err = msgsvc.CreateMessage(message)
	c.Assert(err, Equals, nil)
	msgs, _ := msgsvc.GetMessages(channel, nil, 0, -1)
	c.Assert(len(msgs), Equals, 1)
}

func (s *TestSuite) TestDeleteMessage(c *C) {
	chsvc := s.serviceGroup.ChannelService
	team := NewTeam("superteam", "superorg", "Super Team")
	channel := NewChannel(team, "", "test", "group")
	err := chsvc.SaveChannel(channel, true)

	usersvc := s.serviceGroup.UserService
	msgsvc := s.serviceGroup.MessageService
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
*/
