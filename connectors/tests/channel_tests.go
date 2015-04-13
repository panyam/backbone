package connectors

import (
	"fmt"
	"log"
	// authcore "github.com/panyam/relay/services/auth/core"
	. "github.com/panyam/relay/services/msg/core"
	. "gopkg.in/check.v1"
)

/**
 * Channel creation:
 *
 * Without login - return 4xx
 * With login - success
 */
func (s *TestSuite) TestCreateChannels(c *C) {
	// No login returns Not allowed
	channel, err := s.client.CreateChannel(s.testTeam, "testchannel", true, nil)
	c.Assert(err, Not(IsNil))
	c.Assert(err.Error(), Equals, "401 Unauthorized")
	c.Assert(channel, IsNil)

	// Login and repeat above
	s.LoginClient()
	channel, err = s.client.CreateChannel(s.testTeam, "testchannel", true, []string{"1", "2", "3", "4"})
	c.Assert(err, IsNil)
	c.Assert(channel, Not(IsNil))
	c.Assert(channel.Name, Equals, "testchannel")

	channel.Team = s.testTeam
	channelMembers := s.serviceGroup.ChannelService.GetChannelMembers(channel)
	c.Assert(len(channelMembers), Equals, 5)
}

/**
 * Get channels
 *
 * Public team - return it with or without login
 * Private team:
 * 	Without login - return 4xx
 * 	With login if part of team otherwise 4xx
 */
func (s *TestSuite) TestGetChannels(c *C) {
	team := s.testTeam
	users := make([]*User, 0, 0)
	channels := make([]*Channel, 0, 0)
	for i := 1; i <= 10; i++ {
		creator := NewUser(int64(i), fmt.Sprintf("%d", i), team)
		_ = s.serviceGroup.UserService.SaveUser(creator, false)
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
		members := make([]string, 0, 4)
		for j := 0; j < 5; j++ {
			members = append(members, users[(i+j-1)%len(users)].Username)
		}
		s.serviceGroup.ChannelService.AddChannelMembers(channels[i-1], members)
	}

	// create a channel
	s.LoginClient()
	fetched_channels, channel_members, err := s.client.GetChannels(team, users[0].Username, nil, true, "")
	c.Assert(err, IsNil)
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

/**
 * Get channel details:
 *
 * Public team - return it with or without login
 * Private team:
 * 	Without login - return 4xx
 * 	With login if part of team otherwise 4xx
 */
func (s *TestSuite) TestGetChannelDetails(c *C) {
	// create a channel
	s.LoginClient()
	channel, err := s.client.CreateChannel(s.testTeam, "testchannel", true, []string{"1", "2", "3", "4"})
	c.Assert(err, IsNil)
	c.Assert(channel, Not(IsNil))
	c.Assert(channel.Name, Equals, "testchannel")

	s.LogoutClient()
}

/**
 * Invite to a team
 *
 * 	Without login - return 4xx
 * 	With login if user has "invite" permission and an invitation record is
 * 	created (nothing happens if user is already invited by this or another user)
 */
func (s *TestSuite) TestInviteToChannel(c *C) {
}

/**
 * Join a team
 *
 * 	Without login - return 4xx
 * 	With login:
 * 		if group is public then allow.  Duplicate joins do nothing.
 * 		if group is private and an invitation record exists then allow.
 * 		otherwise 4xx
 */
func (s *TestSuite) TestJoinChannel(c *C) {
}

/**
 * Leave a team
 *
 * 	Without login - return 4xx
 * 	With login:
 * 		if not part of the group - nothing and 200
 * 		if part of the group then leave and 200
 */
func (s *TestSuite) TestLeaveChannel(c *C) {
}
