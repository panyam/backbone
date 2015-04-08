package connectors

import (
	goclient "github.com/panyam/relay/clients/goclient"
	connectorcore "github.com/panyam/relay/connectors/core"
	authcore "github.com/panyam/relay/services/auth/core"
	msgcore "github.com/panyam/relay/services/msg/core"
	. "gopkg.in/check.v1"

	"fmt"
	"testing"
)

type TestSuite struct {
	server       connectorcore.Server
	client       *goclient.ApiClient
	testTeam     *msgcore.Team
	testUser     *msgcore.User
	authService  authcore.IAuthService
	serviceGroup *msgcore.ServiceGroup
	ServerPort   int
	DebugUserId  int64
}

var _ = Suite(&TestSuite{ServerPort: 8000, DebugUserId: 666})

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

func (s *TestSuite) SetUpSuite(c *C) {
	s.serviceGroup, s.authService = s.CreateTestServices()
	s.server = s.CreateTestServer()
	s.client = goclient.NewApiClient(fmt.Sprintf("http://localhost:%d/api", s.ServerPort))
	s.server.SetServiceGroup(s.serviceGroup)
	s.server.SetAuthService(s.authService)
	go s.server.Run()
}

func (s *TestSuite) SetUpTest(c *C) {
	s.serviceGroup.ChannelService.RemoveAllChannels()
	s.serviceGroup.TeamService.RemoveAllTeams()
	s.serviceGroup.UserService.RemoveAllUsers()
	s.serviceGroup.MessageService.RemoveAllMessages()

	s.testTeam, _ = s.serviceGroup.TeamService.CreateTeam(1, "org", "testteam")
	s.testUser = msgcore.NewUser(0, "testuser", s.testTeam)
	err := s.serviceGroup.UserService.SaveUser(s.testUser, false)
	c.Assert(err, IsNil)

	// Also create the debug user for auth
	debugUser := msgcore.NewUser(s.DebugUserId, "dbguser", s.testTeam)
	err = s.serviceGroup.UserService.SaveUser(debugUser, false)
	c.Assert(err, IsNil)

	s.client.DisableAuthentication()
}

func (s *TestSuite) TearDownSuite(c *C) {
	s.server.Stop()
}
