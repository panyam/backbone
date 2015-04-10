package services

import (
	"fmt"
	. "github.com/panyam/relay/services/msg/core"
	. "gopkg.in/check.v1"
	// "code.google.com/p/gomock/gomock"
	"log"
	// "time"
)

func (s *TestSuite) TestSaveChannelNew(c *C) {
	svc := s.serviceGroup.UserService
	team, err := s.serviceGroup.TeamService.CreateTeam(1, "org", "team")
	user := NewUser(0, "user1", team)
	_ = svc.SaveUser(user, false)

	channel := NewChannel(team, user, 0, "test")
	err = s.serviceGroup.ChannelService.SaveChannel(channel, true)

	c.Assert(err, IsNil)
	c.Assert(channel, Not(IsNil))
	c.Assert(channel.Id, Not(Equals), "")
	c.Assert(channel.Name, Equals, "test")
	channel, err = s.serviceGroup.ChannelService.GetChannelById(channel.Id)
	c.Assert(err, IsNil)
	c.Assert(channel, Not(IsNil))
	c.Assert(channel.Name, Equals, "test")
}

func (s *TestSuite) TestSaveChannelExistsById(c *C) {
	svc := s.serviceGroup.UserService
	team, err := s.serviceGroup.TeamService.CreateTeam(1, "org", "team")
	c.Assert(err, IsNil)
	c.Assert(team, Not(IsNil))
	user := NewUser(0, "user1", team)
	_ = svc.SaveUser(user, false)
	channel := NewChannel(team, user, 0, "test")
	_ = s.serviceGroup.ChannelService.SaveChannel(channel, true)

	err = s.serviceGroup.ChannelService.SaveChannel(channel, true)
	c.Assert(err, IsNil)
}

func (s *TestSuite) TestDeleteChannel(c *C) {
	svc := s.serviceGroup.UserService
	team, err := s.serviceGroup.TeamService.CreateTeam(1, "org", "team")
	user := NewUser(0, "user1", team)
	_ = svc.SaveUser(user, false)
	channel := NewChannel(team, user, 0, "test")
	_ = s.serviceGroup.ChannelService.SaveChannel(channel, true)

	c.Assert(err, IsNil)
	c.Assert(channel.Id, Not(Equals), 0)

	err = s.serviceGroup.ChannelService.DeleteChannel(channel)
	c.Assert(err, IsNil)

	channel, err = s.serviceGroup.ChannelService.GetChannelById(channel.Id)
	c.Assert(err, Not(IsNil))
	c.Assert(channel, IsNil)
}

func (s *TestSuite) TestJoinChannel(c *C) {
	svc := s.serviceGroup.UserService
	team, _ := s.serviceGroup.TeamService.CreateTeam(1, "org", "team")
	user := NewUser(0, "user1", team)
	_ = svc.SaveUser(user, false)
	user2 := NewUser(0, "user2", team)
	_ = svc.SaveUser(user2, false)

	channel := NewChannel(team, user, 0, "test")
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
	channel := NewChannel(team, user, 0, "test")
	_ = s.serviceGroup.ChannelService.SaveChannel(channel, true)

	c.Assert(s.serviceGroup.ChannelService.ContainsUser(channel, user), Equals, false)
	s.serviceGroup.ChannelService.JoinChannel(channel, user)
	c.Assert(s.serviceGroup.ChannelService.ContainsUser(channel, user), Equals, true)
}

/**
 * Test searching of channels
 */
func (s *TestSuite) TestGetChannels(c *C) {
	svc := s.serviceGroup.UserService
	team, _ := s.serviceGroup.TeamService.CreateTeam(1, "org", "team")

	users := make([]*User, 0, 0)
	for i := 1; i <= 10; i++ {
		creator := NewUser(int64(i), fmt.Sprintf("user%d", i), team)
		_ = svc.SaveUser(creator, false)
		users = append(users, creator)
		channel := NewChannel(team, creator, int64(i), fmt.Sprintf("channel%d", i))
		err := s.serviceGroup.ChannelService.SaveChannel(channel, true)
		log.Println("Err: ", err)

		// add the creator and 4 next users as members
		members := make([]*User, 0, 4)
		for j := 0; j < 5; j++ {
			members = append(members, users[(i+j)%len(users)])
		}
		s.serviceGroup.ChannelService.AddChannelMembers(channel, members)
	}
}
