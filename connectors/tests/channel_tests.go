package connectors

import (
	// authcore "github.com/panyam/relay/services/auth/core"
	// msgcore "github.com/panyam/relay/services/msg/core"
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
 * Getting teams.  Cases:
 *
 * Without login - return 4xx
 * With login - return teams user is subscribed to or invited to.
 */
func (s *TestSuite) TestGetChannels(c *C) {
}

/**
 * Get team details:
 *
 * Public team - return it with or without login
 * Private team:
 * 	Without login - return 4xx
 * 	With login if part of team otherwise 4xx
 */
func (s *TestSuite) TestGetChannelDetails(c *C) {
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
