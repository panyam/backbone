package connectors

import (
	// . "github.com/panyam/backbone/models"
	. "gopkg.in/check.v1"
	"net/http"
	"testing"
	// "code.google.com/p/gomock/gomock"
	// "log"
	// "time"
)

type TestSuite struct {
	server Server
}

var _ = Suite(&TestSuite{})

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

func (s *TestSuite) SetUpSuite(c *C) {
	s.server = CreateServer()
	go s.server.Run()
}

func (s *TestSuite) TearDownSuite(c *C) {
	s.server.Stop()
}

func (s *TestSuite) SetUpTest(c *C) {
	// clear all items from
}

func (s *TestSuite) TearDownTest(c *C) {
}

func (s *TestSuite) TestUserRegistration(c *C) {
	client := &http.Client{}
	req, err :=
		http.Post("/api/users/register/")
}

func (s *TestSuite) TestUserConfirmation(c *C) {
}
