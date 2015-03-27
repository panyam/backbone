package services

import (
	. "github.com/panyam/backbone/services/core"
	. "gopkg.in/check.v1"
	// "code.google.com/p/gomock/gomock"
	// "log"
	// "time"
)

func (s *TestSuite) TestSaveChannelNew(c *C) {
	svc := s.serviceGroup.UserService
	team, err := s.serviceGroup.TeamService.CreateTeam(1, "org", "team")
	user := NewUser(0, "user1", team)
	_ = svc.SaveUser(user, false)

	channel := NewChannel(team, user, 0, "test", "group")
	err = s.serviceGroup.ChannelService.SaveChannel(channel, true)

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
	svc := s.serviceGroup.UserService
	team, err := s.serviceGroup.TeamService.CreateTeam(1, "org", "team")
	user := NewUser(0, "user1", team)
	_ = svc.SaveUser(user, false)
	channel := NewChannel(team, user, 0, "test", "group")
	_ = s.serviceGroup.ChannelService.SaveChannel(channel, true)

	//err = s.serviceGroup.ChannelService.SaveChannel(channel, false)
	//c.Assert(err, Not(Equals), nil)
	err = s.serviceGroup.ChannelService.SaveChannel(channel, true)
	c.Assert(err, Equals, nil)
}

func (s *TestSuite) TestDeleteChannel(c *C) {
	svc := s.serviceGroup.UserService
	team, err := s.serviceGroup.TeamService.CreateTeam(1, "org", "team")
	user := NewUser(0, "user1", team)
	_ = svc.SaveUser(user, false)
	channel := NewChannel(team, user, 0, "test", "group")
	_ = s.serviceGroup.ChannelService.SaveChannel(channel, true)

	c.Assert(err, Equals, nil)
	c.Assert(channel.Id, Not(Equals), 0)

	err = s.serviceGroup.ChannelService.DeleteChannel(channel)
	c.Assert(err, Equals, nil)

	channel, err = s.serviceGroup.ChannelService.GetChannelById(channel.Id)
	c.Assert(err, Not(Equals), nil)
	c.Assert(channel, Equals, (*Channel)(nil))
}

func (s *TestSuite) TestJoinChannel(c *C) {
	svc := s.serviceGroup.UserService
	team, _ := s.serviceGroup.TeamService.CreateTeam(1, "org", "team")
	user := NewUser(0, "user1", team)
	_ = svc.SaveUser(user, false)
	user2 := NewUser(0, "user2", team)
	_ = svc.SaveUser(user2, false)

	channel := NewChannel(team, user, 0, "test", "group")
	_ = s.serviceGroup.ChannelService.SaveChannel(channel, true)

	s.serviceGroup.ChannelService.JoinChannel(channel, user2)
	c.Assert(channel.ContainsUser(user), Equals, true)
	c.Assert(channel.ContainsUser(user2), Equals, true)
}

/*
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
