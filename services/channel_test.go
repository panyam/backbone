package services

import (
	. "github.com/panyam/backbone/models"
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

func (s *TestSuite) TestCreateChannel(c *C) {
	svc := CreateChannelService()
	channel, err := svc.CreateChannel("", "group", "test")
	c.Assert(err, Equals, nil)
	c.Assert(channel, Not(Equals), (*Channel)(nil))
	c.Assert(channel.Name, Equals, "test")
	channel, err = svc.GetChannelByName("group", "test")
	c.Assert(err, Equals, nil)
	c.Assert(channel, Not(Equals), (*Channel)(nil))
	c.Assert(channel.Name, Equals, "test")
}

func (s *TestSuite) TestCreateChannelExistsByName(c *C) {
	svc := CreateChannelService()
	channel, err := svc.CreateChannel("", "group", "test")
	channel, err = svc.CreateChannel("", "group", "test")
	c.Assert(err, Not(Equals), nil)
	c.Assert(channel, Equals, (*Channel)(nil))
}

func (s *TestSuite) TestCreateChannelExistsById(c *C) {
	svc := CreateChannelService()
	channel, err := svc.CreateChannel("1", "group", "test")
	channel, err = svc.CreateChannel("1", "group2", "test2")
	c.Assert(err, Not(Equals), nil)
	c.Assert(channel, Equals, (*Channel)(nil))
}

func (s *TestSuite) TestDeleteChannel(c *C) {
	svc := CreateChannelService()
	channel, err := svc.CreateChannel("", "group", "test")
	c.Assert(err, Equals, nil)
	c.Assert(channel, Not(Equals), (*Channel)(nil))
	svc.DeleteChannel(channel)
	channel, err = svc.GetChannelByName("group", "test")
	c.Assert(err, Not(Equals), nil)
	c.Assert(channel, Equals, (*Channel)(nil))
}

func (s *TestSuite) TestJoinChannel(c *C) {
	svc := CreateChannelService()
	channel, _ := svc.CreateChannel("", "group", "test")
	user := NewUser("1", "user1")
	svc.JoinChannel(channel, user)
	c.Assert(svc.ChannelContains(channel, user), Equals, true)
}

func (s *TestSuite) TestLeaveChannel(c *C) {
	svc := CreateChannelService()
	channel, _ := svc.CreateChannel("", "group", "test")
	user := NewUser("1", "user1")
	svc.JoinChannel(channel, user)
	c.Assert(svc.ChannelContains(channel, user), Equals, true)
	svc.LeaveChannel(channel, user)
	c.Assert(svc.ChannelContains(channel, user), Equals, false)
}
