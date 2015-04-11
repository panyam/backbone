package services

import (
	"fmt"
	. "github.com/panyam/relay/services/msg/core"
	. "gopkg.in/check.v1"
	"log"
	// "code.google.com/p/gomock/gomock"
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
	channels := make([]*Channel, 0, 0)
	for i := 1; i <= 10; i++ {
		creator := NewUser(int64(i), fmt.Sprintf("%d", i), team)
		_ = svc.SaveUser(creator, false)
		users = append(users, creator)
		channel := NewChannel(team, creator, int64(i), fmt.Sprintf("channel%d", i))
		err := s.serviceGroup.ChannelService.SaveChannel(channel, true)
		if err != nil {
			log.Println("SaveChannel Error: ", err)
		}
		channels = append(channels, channel)
	}

	for i := 1; i <= 10; i++ {
		// add the creator and 4 next users as members
		members := make([]*User, 0, 4)
		for j := 0; j < 5; j++ {
			members = append(members, users[(i+j-1)%len(users)])
		}
		s.serviceGroup.ChannelService.AddChannelMembers(channels[i-1], members)
	}

	// Test owner filter
	fetched_channels, channel_members := s.serviceGroup.ChannelService.GetChannels(team, users[0], "", nil, true)
	c.Assert(len(fetched_channels), Equals, 1)
	c.Assert(len(channel_members), Equals, 1)
	// ensure all users have the same creator
	c.Assert(fetched_channels[0].Creator.Id, Equals, users[0].Id)
	c.Assert(len(channel_members[0]), Equals, 5)

	// Test participants
	fetched_channels, channel_members = s.serviceGroup.ChannelService.GetChannels(team, nil, "", []*User{users[1], users[2]}, true)
	c.Assert(len(fetched_channels), Equals, 4)
	for i := 0; i < 4; i++ {
		c.Assert(len(channel_members[i]), Equals, 5)
	}
}
