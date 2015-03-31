package connectors

import (
	service_core "github.com/panyam/relay/connectors/core"
	// client "github.com/panyam/relay/services/client/goclient"
	connector_core "github.com/panyam/relay/services/core"
	. "gopkg.in/check.v1"
	"testing"
)

type TestSuite struct {
	server       service_core.Server
	serviceGroup *connector_core.ServiceGroup
}

var _ = Suite(&TestSuite{})

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

func (s *TestSuite) SetUpSuite(c *C) {
	s.serviceGroup = CreateTestServiceGroup()
	s.server = CreateTestServer()
	s.server.SetServiceGroup(s.serviceGroup)
	go s.server.Run()
}

func (s *TestSuite) SetUpTest(c *C) {
	s.serviceGroup.ChannelService.RemoveAllChannels()
	s.serviceGroup.TeamService.RemoveAllTeams()
	s.serviceGroup.UserService.RemoveAllUsers()
	s.serviceGroup.MessageService.RemoveAllMessages()
}

func (s *TestSuite) TearDownSuite(c *C) {
	s.server.Stop()
}

func (s *TestSuite) TestUserRegistration(c *C) {
}
