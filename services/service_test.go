package services

import (
	// . "github.com/panyam/backbone/models"
	. "gopkg.in/check.v1"
	"testing"
	// "code.google.com/p/gomock/gomock"
	// "log"
	// "time"
)

type TestSuite struct{}

var _ = Suite(&TestSuite{})

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }
