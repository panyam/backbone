package services

/*
import (
	. "github.com/panyam/backbone/services/core"
	. "gopkg.in/check.v1"
	// "code.google.com/p/gomock/gomock"
	// "log"
	// "time"
)

func (s *TestSuite) TestSaveChannelNew(c *C) {
	team := NewTeam("superteam", "superorg", "Super Team")
	channel := NewChannel(team, "", "test", "group")
	err := s.serviceGroup.ChannelService.SaveChannel(channel, true)
	c.Assert(err, Equals, nil)
	c.Assert(channel, Not(Equals), (*Channel)(nil))
	c.Assert(channel.Id, Not(Equals), "")
	c.Assert(channel.Name, Equals, "test")
	channel, err = s.serviceGroup.ChannelService.GetChannelById(channel.Id)
	c.Assert(err, Equals, nil)
	c.Assert(channel, Not(Equals), (*Channel)(nil))
	c.Assert(channel.Name, Equals, "test")
}

func (s *TestSuite) TestSaveChannelExistsById(c *C) {
	team := NewTeam("superteam", "superorg", "Super Team")
	channel := NewChannel(team, "", "test", "group")
	err := s.serviceGroup.ChannelService.SaveChannel(channel, true)
	c.Assert(err, Equals, nil)
	err = s.serviceGroup.ChannelService.SaveChannel(channel, false)
	c.Assert(err, Not(Equals), nil)
	err = s.serviceGroup.ChannelService.SaveChannel(channel, true)
	c.Assert(err, Equals, nil)
}

func (s *TestSuite) TestDeleteChannel(c *C) {
	team := NewTeam("superteam", "superorg", "Super Team")
	channel := NewChannel(team, "", "test", "group")
	err := s.serviceGroup.ChannelService.SaveChannel(channel, true)

	c.Assert(err, Equals, nil)
	c.Assert(channel, Not(Equals), (*Channel)(nil))
	s.serviceGroup.ChannelService.DeleteChannel(channel)
	channel, err = s.serviceGroup.ChannelService.GetChannelById(channel.Id)
	c.Assert(err, Not(Equals), nil)
	c.Assert(channel, Equals, (*Channel)(nil))
}

func (s *TestSuite) TestJoinChannel(c *C) {
	team := NewTeam("superteam", "superorg", "Super Team")
	channel := NewChannel(team, "", "test", "group")
	s.serviceGroup.ChannelService.SaveChannel(channel, true)

	user := NewUser("1", "user1")
	s.serviceGroup.ChannelService.JoinChannel(channel, user)
	c.Assert(channel.ContainsUser(user), Equals, true)
}

func (s *TestSuite) TestLeaveChannel(c *C) {
	team := NewTeam("superteam", "superorg", "Super Team")
	channel := NewChannel(team, "", "test", "group")
	s.serviceGroup.ChannelService.SaveChannel(channel, true)

	user := NewUser("1", "user1")
	s.serviceGroup.ChannelService.JoinChannel(channel, user)
	c.Assert(channel.ContainsUser(user), Equals, true)
	s.serviceGroup.ChannelService.LeaveChannel(channel, user)
	c.Assert(channel.ContainsUser(user), Equals, false)
}
*/
