package services

import (
	"github.com/panyam/relay/services/msg/core"
	. "gopkg.in/check.v1"
	"testing"
)

type TestSuite struct {
	serviceGroup core.ServiceGroup
}

var _ = Suite(&TestSuite{})

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

func (s *TestSuite) SetUpSuite(c *C) {
	s.serviceGroup = CreateServiceGroup()
}

func (s *TestSuite) SetUpTest(c *C) {
	s.serviceGroup.IDService.RemoveAllDomains(nil)
	s.serviceGroup.ChannelService.RemoveAllChannels(nil)
	s.serviceGroup.TeamService.RemoveAllTeams(nil)
	s.serviceGroup.UserService.RemoveAllUsers(nil)
	s.serviceGroup.MessageService.RemoveAllMessages(nil)
}

func (s *TestSuite) TearDownTest(c *C) {
}
