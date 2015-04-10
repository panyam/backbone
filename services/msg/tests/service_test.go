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
	s.serviceGroup.IDService.RemoveAllDomains()
	s.serviceGroup.ChannelService.RemoveAllChannels()
	s.serviceGroup.TeamService.RemoveAllTeams()
	s.serviceGroup.UserService.RemoveAllUsers()
	s.serviceGroup.MessageService.RemoveAllMessages()
}

func (s *TestSuite) TearDownTest(c *C) {
}
