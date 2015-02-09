package service

import (
	// "code.google.com/p/gomock/gomock"
	// . "github.com/panyam/backbone/models"
	. "gopkg.in/check.v1"
	"testing"
	// "time"
)

type ServiceSuite struct{}

var _ = Suite(&ServiceSuite{})

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

func (s *ServiceSuite) TestGetTeams(c *C) {
	c.Assert(0, Equals, 1)
}

// How to write these damn tests?
// User centric POV may help?  Assume user exists.  This is all about messaging.
// So start with that instead of teams/channels etc
// How to test and how/when to care about persistence?   Service already has the
// operations that need to be implemented.  Where the hell do we start.
// Gaaaaaah.

func (s *ServiceSuite) TestSendMessage(c *C) {
	// u1 := NewUser()
}
