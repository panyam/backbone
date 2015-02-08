package service

import (
	// "code.google.com/p/gomock/gomock"
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
