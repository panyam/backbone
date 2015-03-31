package services

import (
	. "github.com/panyam/relay/services/core"
	. "gopkg.in/check.v1"
	// "code.google.com/p/gomock/gomock"
	"fmt"
	"log"
	// "time"
)

func (s *TestSuite) TestCreateMessageService(c *C) {
	svc := s.serviceGroup.MessageService
	c.Assert(svc, Not(Equals), nil)
}

func (s *TestSuite) MakeTestChannel() *Channel {
	// create team
	team, _ := s.serviceGroup.TeamService.CreateTeam(1, "org", "team")

	user := NewUser(0, "user1", team)
	_ = s.serviceGroup.UserService.SaveUser(user, false)

	// create channel
	channel := NewChannel(team, user, 0, "test", "group")
	s.serviceGroup.ChannelService.SaveChannel(channel, true)

	return channel
}

func (s *TestSuite) TestGetMessages(c *C) {
	channel := s.MakeTestChannel()

	msgs, _ := s.serviceGroup.MessageService.GetMessages(channel, nil, 0, -1)
	c.Assert(len(msgs), Equals, 0)
}

func (s *TestSuite) TestSaveMessage(c *C) {
	channel := s.MakeTestChannel()

	message := NewMessage(channel, channel.Creator)
	err := s.serviceGroup.MessageService.SaveMessage(message)
	c.Assert(err, Equals, nil)
	msgs, err := s.serviceGroup.MessageService.GetMessages(channel, nil, 0, -1)
	log.Println("err: ", err)
	c.Assert(len(msgs), Equals, 1)
}

/**
 * Test pagination
 */
func (s *TestSuite) TestPaginationAtFront(c *C) {
	channel := s.MakeTestChannel()

	for i := 0; i < 100; i++ {
		message := NewMessage(channel, channel.Creator)
		message.MsgData = fmt.Sprintf("Message %d", i)
		err := s.serviceGroup.MessageService.SaveMessage(message)
		c.Assert(err, Equals, nil)
	}

	top20, err := s.serviceGroup.MessageService.GetMessages(channel, nil, 0, 20)
	c.Assert(err, Equals, nil)
	c.Assert(len(top20), Equals, 20)

	mid20, err := s.serviceGroup.MessageService.GetMessages(channel, nil, 30, 20)
	c.Assert(err, Equals, nil)
	c.Assert(len(mid20), Equals, 20)

	bottom20, err := s.serviceGroup.MessageService.GetMessages(channel, nil, -1, 20)
	c.Assert(err, Equals, nil)
	c.Assert(len(bottom20), Equals, 20)

	// check messages are correct
	for i := 0; i < 20; i++ {
		c.Assert(top20[i].MsgData, Equals, fmt.Sprintf("Message %d", i))
	}

	for i := 30; i < 50; i++ {
		c.Assert(mid20[i-30].MsgData, Equals, fmt.Sprintf("Message %d", i))
	}

	for i := 80; i < 100; i++ {
		c.Assert(bottom20[i-80].MsgData, Equals, fmt.Sprintf("Message %d", i))
	}
}

func (s *TestSuite) TestDeleteMessage(c *C) {
	channel := s.MakeTestChannel()

	message := NewMessage(channel, channel.Creator)
	err := s.serviceGroup.MessageService.SaveMessage(message)
	c.Assert(err, Equals, nil)

	msgs, err := s.serviceGroup.MessageService.GetMessages(channel, nil, 0, -1)
	c.Assert(len(msgs), Equals, 1)

	err = s.serviceGroup.MessageService.DeleteMessage(message)
	c.Assert(err, Equals, nil)
	msgs, err = s.serviceGroup.MessageService.GetMessages(channel, nil, 0, -1)
	c.Assert(len(msgs), Equals, 0)
}
