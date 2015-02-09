package client

import (
	// "code.google.com/p/gomock/gomock"
	"github.com/panyam/backbone/core"
	. "gopkg.in/check.v1"
	"testing"
	// "time"
)

type ClientSuite struct{}

var _ = Suite(&ClientSuite{})

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

func CreateNewClient() Client {
	return new(NativeClient)
}

func LoggedInClient(username string) Client {
	client := CreateNewClient()
	client.Connect()
	if username == "" {
		client.Login(nil)
	} else {
		client.Login(map[string]string{"username": username})
	}
	return client
}

func (s *ClientSuite) TestCreateClient(c *C) {
	client := CreateNewClient()
	c.Assert(client, Not(Equals), nil)
	c.Assert(client.IsConnected(), Equals, false)
	c.Assert(client.CurrentUser(), Equals, nil)
}

func (s *ClientSuite) TestConnection(c *C) {
	client := CreateNewClient()
	error := client.Connect()
	c.Assert(client.IsConnected(), Equals, true)
	c.Assert(client.CurrentUser(), Equals, nil)
	c.Assert(error, Equals, nil)
}

func (s *ClientSuite) TestTeams(c *C) {
	client := CreateNewClient()
	client.Connect()
	teams, _ := client.Teams()
	c.Assert(len(teams), Equals, 0)
}

func (s *ClientSuite) TestLoginWithNoCredentials(c *C) {
	client := LoggedInClient("")
	c.Assert(client.CurrentUser(), Equals, nil)
}

func (s *ClientSuite) TestLoginSuccess(c *C) {
	client := LoggedInClient("sri")
	c.Assert(client.CurrentUser(), Not(Equals), core.User(nil))
	c.Assert(client.CurrentUser().Username(), Equals, "sri")
}

func (s *ClientSuite) TestLogout(c *C) {
	client := LoggedInClient("sri")
	c.Assert(client.CurrentUser(), Not(Equals), core.User(nil))
	client.Logout()
	c.Assert(client.CurrentUser(), Equals, nil)
}

// Once a user is logged in, he should be able to see all his teams
// First time the user should have no teams (this requires dynamic user creation
// and deletion)
func (s *ClientSuite) TestTeamsForNewUser(c *C) {
	client := LoggedInClient("newuser")
	teams, error := client.Teams()
	c.Assert(error, Equals, nil)
	c.Assert(len(teams), Equals, 0)
}

// A user who is admin can create a team
// Requires: admin user with no teams
func (s *ClientSuite) TestTeamCreationByAdminUser(c *C) {
}
