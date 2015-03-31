package connectors

import (
	. "gopkg.in/check.v1"
)

/**
 * Test that after a user registration, we
 */
func (s *TestSuite) TestUserRegistration(c *C) {
	s.client.RegisterUser("testuser", "phone", "1231231234", "password")
}
