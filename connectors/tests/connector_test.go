package connectors

import (
	goclient "github.com/panyam/relay/clients/goclient"
	connector_core "github.com/panyam/relay/connectors/core"
	auth_core "github.com/panyam/relay/services/auth/core"
	msg_core "github.com/panyam/relay/services/msg/core"
	. "gopkg.in/check.v1"
	"testing"
)

type TestSuite struct {
	server       connector_core.Server
	client       *goclient.ApiClient
	authService  auth_core.IAuthService
	serviceGroup *msg_core.ServiceGroup
}

var _ = Suite(&TestSuite{})

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

func (s *TestSuite) SetUpSuite(c *C) {
	s.serviceGroup, s.authService = CreateTestServices()
	s.server = CreateTestServer()
	s.client = goclient.NewApiClient("http://localhost:3000/api")
	s.client.Authenticator = &goclient.DebugAuthenticator{"testuser"}
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
