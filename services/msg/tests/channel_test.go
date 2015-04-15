package services

import (
	"fmt"
	msgcore "github.com/panyam/relay/services/msg/core"
	. "gopkg.in/check.v1"
	"log"
	// "code.google.com/p/gomock/gomock"
	// "time"
)

func (s *TestSuite) TestCreateChannelNew(c *C) {
	svc := s.serviceGroup.UserService
	team, err := s.serviceGroup.TeamService.SaveTeam(msgcore.NewTeam(1, "org", "team"))
	user := msgcore.NewUserByName("user1", team)
	_ = svc.SaveUser(&msgcore.SaveUserRequest{nil, user, false})

	channel := msgcore.NewChannel(team, user, 0, "test")
	channel, err = s.serviceGroup.ChannelService.CreateChannel(&msgcore.CreateChannelRequest{channel, true})

	c.Assert(err, IsNil)
	c.Assert(channel, Not(IsNil))
	c.Assert(channel.Id, Not(Equals), "")
	c.Assert(channel.Name, Equals, "test")
	result, err := s.serviceGroup.ChannelService.GetChannelById(channel.Id)
	c.Assert(err, IsNil)
	c.Assert(result, Not(IsNil))
	c.Assert(result.Channel, Not(IsNil))
	c.Assert(result.Channel.Name, Equals, "test")
}

func (s *TestSuite) TestCreateChannelExistsById(c *C) {
	svc := s.serviceGroup.UserService
	team, err := s.serviceGroup.TeamService.SaveTeam(msgcore.NewTeam(1, "org", "team"))
	c.Assert(err, IsNil)
	c.Assert(team, Not(IsNil))
	user := msgcore.NewUserByName("user1", team)
	_ = svc.SaveUser(&msgcore.SaveUserRequest{nil, user, false})
	channel := msgcore.NewChannel(team, user, 0, "test")
	channel, _ = s.serviceGroup.ChannelService.CreateChannel(&msgcore.CreateChannelRequest{channel, true})

	channel, err = s.serviceGroup.ChannelService.CreateChannel(&msgcore.CreateChannelRequest{channel, true})
	c.Assert(err, IsNil)
}

func (s *TestSuite) TestDeleteChannel(c *C) {
	svc := s.serviceGroup.UserService
	team, err := s.serviceGroup.TeamService.SaveTeam(msgcore.NewTeam(1, "org", "team"))
	user := msgcore.NewUserByName("user1", team)
	_ = svc.SaveUser(&msgcore.SaveUserRequest{nil, user, false})
	channel := msgcore.NewChannel(team, user, 0, "test")
	channel, _ = s.serviceGroup.ChannelService.CreateChannel(&msgcore.CreateChannelRequest{channel, true})

	c.Assert(err, IsNil)
	c.Assert(channel.Id, Not(Equals), 0)

	err = s.serviceGroup.ChannelService.DeleteChannel(&msgcore.DeleteChannelRequest{nil, channel})
	c.Assert(err, IsNil)

	result, err := s.serviceGroup.ChannelService.GetChannelById(channel.Id)
	c.Assert(err, Not(IsNil))
	c.Assert(result, IsNil)
}

func (s *TestSuite) TestJoinChannel(c *C) {
	svc := s.serviceGroup.UserService
	team, _ := s.serviceGroup.TeamService.SaveTeam(msgcore.NewTeam(1, "org", "team"))
	user := msgcore.NewUserByName("user1", team)
	_ = svc.SaveUser(&msgcore.SaveUserRequest{nil, user, false})
	user2 := msgcore.NewUserByName("user2", team)
	_ = svc.SaveUser(&msgcore.SaveUserRequest{nil, user2, false})

	channel := msgcore.NewChannel(team, user, 0, "test")
	channel, _ = s.serviceGroup.ChannelService.CreateChannel(&msgcore.CreateChannelRequest{channel, true})

	s.serviceGroup.ChannelService.JoinChannel(&msgcore.ChannelMembershipRequest{nil, channel, user})
	s.serviceGroup.ChannelService.JoinChannel(&msgcore.ChannelMembershipRequest{nil, channel, user2})
	c.Assert(s.serviceGroup.ChannelService.ContainsUser(&msgcore.ChannelMembershipRequest{nil, channel, user2}), Equals, true)
	c.Assert(s.serviceGroup.ChannelService.ContainsUser(&msgcore.ChannelMembershipRequest{nil, channel, user2}), Equals, true)
}

func (s *TestSuite) TestLeaveChannel(c *C) {
	svc := s.serviceGroup.UserService
	team, _ := s.serviceGroup.TeamService.SaveTeam(msgcore.NewTeam(1, "org", "team"))
	user := msgcore.NewUserByName("user1", team)
	_ = svc.SaveUser(&msgcore.SaveUserRequest{nil, user, false})
	channel := msgcore.NewChannel(team, user, 0, "test")
	channel, _ = s.serviceGroup.ChannelService.CreateChannel(&msgcore.CreateChannelRequest{channel, true})

	c.Assert(s.serviceGroup.ChannelService.ContainsUser(&msgcore.ChannelMembershipRequest{nil, channel, user}), Equals, false)
	s.serviceGroup.ChannelService.JoinChannel(&msgcore.ChannelMembershipRequest{nil, channel, user})
	c.Assert(s.serviceGroup.ChannelService.ContainsUser(&msgcore.ChannelMembershipRequest{nil, channel, user}), Equals, true)
}

/**
 * Test searching of channels
 */
func (s *TestSuite) TestGetChannels(c *C) {
	svc := s.serviceGroup.UserService
	team, _ := s.serviceGroup.TeamService.SaveTeam(msgcore.NewTeam(1, "org", "team"))

	users := make([]*msgcore.User, 0, 0)
	channels := make([]*msgcore.Channel, 0, 0)
	for i := 1; i <= 10; i++ {
		creator := msgcore.NewUser(int64(i), fmt.Sprintf("%d", i), team)
		_ = svc.SaveUser(&msgcore.SaveUserRequest{nil, creator, false})
		users = append(users, creator)
		channel := msgcore.NewChannel(team, creator, int64(i), fmt.Sprintf("channel%d", i))
		channel, err := s.serviceGroup.ChannelService.CreateChannel(&msgcore.CreateChannelRequest{channel, true})
		if err != nil {
			log.Println("CreateChannel Error: ", err)
		}
		channels = append(channels, channel)
	}

	for i := 1; i <= 10; i++ {
		// add the creator and 4 next users as members
		members := make([]string, 0, 4)
		for j := 0; j < 5; j++ {
			members = append(members, users[(i+j-1)%len(users)].Username)
		}
		s.serviceGroup.ChannelService.AddChannelMembers(&msgcore.InviteMembersRequest{nil, channels[i-1], members})
	}

	// Test owner filter
	request := &msgcore.GetChannelsRequest{team, users[0], "", nil, true}
	result, _ := s.serviceGroup.ChannelService.GetChannels(request)
	c.Assert(len(result.Channels), Equals, 1)
	c.Assert(len(result.Members), Equals, 1)
	// ensure all users have the same creator
	c.Assert(result.Channels[0].Creator.Id, Equals, users[0].Id)
	c.Assert(len(result.Members[0]), Equals, 5)

	// Test participants
	request = &msgcore.GetChannelsRequest{team, nil, "", []*msgcore.User{users[1], users[2]}, true}
	result, _ = s.serviceGroup.ChannelService.GetChannels(request)
	c.Assert(len(result.Channels), Equals, 4)
	for i := 0; i < 4; i++ {
		c.Assert(len(result.Members[i]), Equals, 5)
	}
}
