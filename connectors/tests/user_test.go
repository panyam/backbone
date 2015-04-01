package connectors

import (
	authcore "github.com/panyam/relay/services/auth/core"
	. "gopkg.in/check.v1"
)

/**
 * Test that after a user registration, we
 */
func (s *TestSuite) TestUserRegistration(c *C) {
	registration, err := s.client.RegisterUser(1, "testuser", "phone", "1231231234", "password")
	c.Assert(err, Equals, nil)
	c.Assert(registration, Not(Equals), (*authcore.Registration)(nil))
	c.Assert(registration.Id, Not(Equals), 0)
}
