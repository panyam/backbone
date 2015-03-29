package services

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

func (s *TestSuite) MakeTestChannel() *Channel {
	// create team
	team, _ := s.serviceGroup.TeamService.CreateTeam(1, "org", "team")

	user := NewUser(0, "user1", team)
	_ = s.serviceGroup.UserService.SaveUser(user, false)

	// create channel
	channel := NewChannel(team, user, 0, "test", "group")
	s.serviceGroup.ChannelService.SaveChannel(channel, true)

	return channel
}

func (s *TestSuite) TestGetMessages(c *C) {
	channel := s.MakeTestChannel()

	msgs, _ := s.serviceGroup.MessageService.GetMessages(channel, nil, 0, -1)
	c.Assert(len(msgs), Equals, 0)
}

func (s *TestSuite) TestCreateMessage(c *C) {
	channel := s.MakeTestChannel()

	message := NewMessage(channel, channel.Creator)
	err := s.serviceGroup.MessageService.CreateMessage(message)
	c.Assert(err, Equals, nil)
	msgs, _ := s.serviceGroup.MessageService.GetMessages(channel, nil, 0, -1)
	c.Assert(len(msgs), Equals, 1)
}

/*
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
