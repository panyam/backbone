package services

import (
	// . "github.com/panyam/backbone/models"
	. "gopkg.in/check.v1"
	// "code.google.com/p/gomock/gomock"
	// "log"
	// "time"
)

func CreateChannelService() IChannelService {
	return NewChannelService()
}

func (s *TestSuite) TestCreateChannelService(c *C) {
	svc := CreateChannelService()
	c.Assert(svc, Not(Equals), nil)
}
