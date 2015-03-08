package services

import (
	. "gopkg.in/check.v1"
	"testing"
)

type TestSuite struct{}

var _ = Suite(&TestSuite{})

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }
