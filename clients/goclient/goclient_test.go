package goclient

import (
	"github.com/panyam/relay/services/core"
	. "gopkg.in/check.v1"
	"testing"
	// "code.google.com/p/gomock/gomock"
)

type TestSuite struct {
	client *ApiClient
	sg     *core.ServiceGroup
}

var _ = Suite(&TestSuite{})

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

func (s *TestSuite) SetUpSuite(c *C)    {}
func (s *TestSuite) TearDownSuite(c *C) {}

func (s *TestSuite) SetUpTest(c *C) {
	s.client = NewApiClient("http://localhost:3000/api")
	s.client.Authenticator = &DebugAuthenticator{"testuser"}
}

func (s *TestSuite) TearDownTest(c *C) {
}

func (s *TestSuite) TestUserRegistration(c *C) {
	s.client.RegisterUser("testuser", "phone", "1231231234", "password")
}
