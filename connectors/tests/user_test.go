package connectors

import (
	authcore "github.com/panyam/relay/services/auth/core"
	msgcore "github.com/panyam/relay/services/msg/core"
	. "gopkg.in/check.v1"
)

/**
 * Test that after a user registration, we
 */
func (s *TestSuite) TestUserRegistration(c *C) {
	registration, err := s.client.RegisterUser(s.testTeam, "user2", "phone", "1231231234", "password")
	c.Assert(err, Equals, nil)
	c.Assert(registration, Not(Equals), (*authcore.Registration)(nil))
	c.Assert(registration.Id, Not(Equals), 0)

	// ensure we can get this registration
	oldId := registration.Id
	registration, err = s.authService.GetRegistrationById(registration.Id)
	c.Assert(err, Equals, nil)
	c.Assert(registration, Not(Equals), (*authcore.Registration)(nil))
	c.Assert(registration.Id, Equals, oldId)

	// ensure that this user is not created - this should happen after
	// confirmation
	user, err := s.serviceGroup.UserService.GetUser("user2", s.testTeam)
	c.Assert(user, Equals, (*msgcore.User)(nil))
}

/**
 * Test confirmation of a registration.
 */
func (s *TestSuite) TestUserConfirmationFailed(c *C) {
	/*
		registration, err := s.client.RegisterUser(1, "testuser", "phone", "1231231234", "password")
		c.Assert(err, Equals, nil)
		c.Assert(registration, Not(Equals), (*authcore.Registration)(nil))
		c.Assert(registration.Id, Not(Equals), 0)

		err = s.client.ConfirmRegistration(registration.Id, "")
		c.Assert(err, Not(Equals), nil)
		c.Assert(err.Error(), Equals, "Confirmation failed")
		user, err := s.serviceGroup.UserService.GetUser("testuser", s.testTeam)
		c.Assert(user, Equals, (*msgcore.User)(nil))
	*/
}
