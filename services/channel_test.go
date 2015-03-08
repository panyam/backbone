package services

import (
	. "github.com/panyam/backbone/models"
	. "gopkg.in/check.v1"
	// "code.google.com/p/gomock/gomock"
	// "log"
	// "time"
)

func (s *TestSuite) TestCreateChannelService(c *C) {
	svc := CreateChannelService()
	c.Assert(svc, Not(Equals), nil)
}

func (s *TestSuite) TestSaveChannelNew(c *C) {
	svc := CreateChannelService()
	team := NewTeam("superteam", "superorg", "Super Team")
	channel := NewChannel(team, "", "test", "group")
	err := svc.SaveChannel(channel, true)
	c.Assert(err, Equals, nil)
	c.Assert(channel, Not(Equals), (*Channel)(nil))
	c.Assert(channel.Id, Not(Equals), "")
	c.Assert(channel.Name, Equals, "test")
	channel, err = svc.GetChannelById(channel.Id)
	c.Assert(err, Equals, nil)
	c.Assert(channel, Not(Equals), (*Channel)(nil))
	c.Assert(channel.Name, Equals, "test")
}

func (s *TestSuite) TestSaveChannelExistsById(c *C) {
	svc := CreateChannelService()
	team := NewTeam("superteam", "superorg", "Super Team")
	channel := NewChannel(team, "", "test", "group")
	err := svc.SaveChannel(channel, true)
	c.Assert(err, Equals, nil)
	err = svc.SaveChannel(channel, false)
	c.Assert(err, Not(Equals), nil)
	err = svc.SaveChannel(channel, true)
	c.Assert(err, Equals, nil)
}

func (s *TestSuite) TestDeleteChannel(c *C) {
	svc := CreateChannelService()
	team := NewTeam("superteam", "superorg", "Super Team")
	channel := NewChannel(team, "", "test", "group")
	err := svc.SaveChannel(channel, true)

	c.Assert(err, Equals, nil)
	c.Assert(channel, Not(Equals), (*Channel)(nil))
	svc.DeleteChannel(channel)
	channel, err = svc.GetChannelById(channel.Id)
	c.Assert(err, Not(Equals), nil)
	c.Assert(channel, Equals, (*Channel)(nil))
}

func (s *TestSuite) TestJoinChannel(c *C) {
	svc := CreateChannelService()
	team := NewTeam("superteam", "superorg", "Super Team")
	channel := NewChannel(team, "", "test", "group")
	svc.SaveChannel(channel, true)

	user := NewUser("1", "user1")
	svc.JoinChannel(channel, user)
	c.Assert(channel.ContainsUser(user), Equals, true)
}

func (s *TestSuite) TestLeaveChannel(c *C) {
	svc := CreateChannelService()
	team := NewTeam("superteam", "superorg", "Super Team")
	channel := NewChannel(team, "", "test", "group")
	svc.SaveChannel(channel, true)

	user := NewUser("1", "user1")
	svc.JoinChannel(channel, user)
	c.Assert(channel.ContainsUser(user), Equals, true)
	svc.LeaveChannel(channel, user)
	c.Assert(channel.ContainsUser(user), Equals, false)
}
