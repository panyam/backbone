package services

import (
	. "github.com/panyam/relay/services/msg/core"
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

	s.serviceGroup.ChannelService.JoinChannel(channel, user)
	s.serviceGroup.ChannelService.JoinChannel(channel, user2)
	c.Assert(s.serviceGroup.ChannelService.ContainsUser(channel, user), Equals, true)
	c.Assert(s.serviceGroup.ChannelService.ContainsUser(channel, user2), Equals, true)
}

func (s *TestSuite) TestLeaveChannel(c *C) {
	svc := s.serviceGroup.UserService
	team, _ := s.serviceGroup.TeamService.CreateTeam(1, "org", "team")
	user := NewUser(0, "user1", team)
	_ = svc.SaveUser(user, false)
	channel := NewChannel(team, user, 0, "test", "group")
	_ = s.serviceGroup.ChannelService.SaveChannel(channel, true)

	c.Assert(s.serviceGroup.ChannelService.ContainsUser(channel, user), Equals, false)
	s.serviceGroup.ChannelService.JoinChannel(channel, user)
	c.Assert(s.serviceGroup.ChannelService.ContainsUser(channel, user), Equals, true)
}