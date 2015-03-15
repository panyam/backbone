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
	svg := CreateServiceGroup()
	svc := svg.MessageService
	c.Assert(svc, Not(Equals), nil)
}

func (s *TestSuite) TestGetMessages(c *C) {
	svg := CreateServiceGroup()
	chsvc := svg.ChannelService
	team := NewTeam("superteam", "superorg", "Super Team")
	channel := NewChannel(team, "", "test", "group")
	chsvc.SaveChannel(channel, true)

	msgsvc := svg.MessageService
	msgs, _ := msgsvc.GetMessages(channel, nil, 0, -1)
	c.Assert(len(msgs), Equals, 0)
}

func (s *TestSuite) TestCreateMessage(c *C) {
	svg := CreateServiceGroup()
	chsvc := svg.ChannelService
	team, _ := svg.TeamService.CreateTeam("1", "org", "team")
	channel := NewChannel(team, "", "test", "group")
	err := chsvc.SaveChannel(channel, true)

	usersvc := svg.UserService
	msgsvc := svg.MessageService

	sender := User{Username: "user1", Team: team}
	err := usersvc.SaveUser(&sender, false)

	message := NewMessage(channel, sender)
	err = msgsvc.CreateMessage(message)
	c.Assert(err, Equals, nil)
	msgs, _ := msgsvc.GetMessages(channel, nil, 0, -1)
	c.Assert(len(msgs), Equals, 1)
}

func (s *TestSuite) TestDeleteMessage(c *C) {
	svg := CreateServiceGroup()
	chsvc := svg.ChannelService
	team := NewTeam("superteam", "superorg", "Super Team")
	channel := NewChannel(team, "", "test", "group")
	err := chsvc.SaveChannel(channel, true)

	usersvc := svg.UserService
	msgsvc := svg.MessageService
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
