package connectors

import (
	// authcore "github.com/panyam/relay/services/auth/core"
	// msgcore "github.com/panyam/relay/services/msg/core"
	. "gopkg.in/check.v1"
)

/**
 * Team creation:
 *
 * Without login - return 4xx
 * With login
 * 	If has "createteam" permission then allow
 * 	Otherwise disallow
 *
 */
func (s *TestSuite) TestCreateTeams(c *C) {
}

/**
 * Getting teams.  Cases:
 *
 * Without login - return 4xx
 * With login - return teams user is subscribed to or invited to.
 */
func (s *TestSuite) TestGetTeams(c *C) {
}

/**
 * Get team details:
 *
 * Public team - return it with or without login
 * Private team:
 * 	Without login - return 4xx
 * 	With login if part of team otherwise 4xx
 */
func (s *TestSuite) TestGetTeamDetails(c *C) {
}

/**
 * Invite to a team
 *
 * 	Without login - return 4xx
 * 	With login if user has "invite" permission and an invitation record is
 * 	created (nothing happens if user is already invited by this or another user)
 */
func (s *TestSuite) TestInviteToTeam(c *C) {
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
func (s *TestSuite) TestJoinTeam(c *C) {
}

/**
 * Leave a team
 *
 * 	Without login - return 4xx
 * 	With login:
 * 		if not part of the group - nothing and 200
 * 		if part of the group then leave and 200
 */
func (s *TestSuite) TestLeaveTeam(c *C) {
}
