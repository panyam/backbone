package connectors

import (
	goclient "github.com/panyam/relay/clients/goclient"
	connectorcore "github.com/panyam/relay/connectors/core"
	authcore "github.com/panyam/relay/services/auth/core"
	msgcore "github.com/panyam/relay/services/msg/core"
	. "gopkg.in/check.v1"
	"testing"
)

type TestSuite struct {
	server       connectorcore.Server
	client       *goclient.ApiClient
	testTeam     *msgcore.Team
	testUser     *msgcore.User
	authService  authcore.IAuthService
	serviceGroup *msgcore.ServiceGroup
}

var _ = Suite(&TestSuite{})

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

func (s *TestSuite) SetUpSuite(c *C) {
	s.serviceGroup, s.authService = s.CreateTestServices()
	s.server = s.CreateTestServer()
	s.client = goclient.NewApiClient("http://localhost:3000/api")
	s.client.Authenticator = &goclient.DebugAuthenticator{"testuser"}
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
	c.Assert(err, Equals, nil)

	s.client.DisableAuthentication()
}

func (s *TestSuite) TearDownSuite(c *C) {
	s.server.Stop()
}
