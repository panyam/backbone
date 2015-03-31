package inmem

import (
	"errors"
	. "github.com/panyam/relay/services/messaging/core"
)

type ChannelService struct {
	Cls            interface{}
	channelCounter int64
	channelsById   map[int64]*Channel
}

func NewChannelService() *ChannelService {
	svc := ChannelService{}
	svc.Cls = &svc
	svc.RemoveAllChannels()
	return &svc
}

/**
 * Lets a user create a channel.
 */
func (c *ChannelService) SaveChannel(channel *Channel, override bool) error {
	if channel.Id == 0 {
		channel.Id = c.channelCounter
		c.channelCounter++
	} else if _, ok := c.channelsById[channel.Id]; ok && !override {
		return errors.New("Channel already exists by ID")
	}
	c.channelsById[channel.Id] = channel
	return nil
}

/**
 * Retrieve a channel by Name.
 */
func (c *ChannelService) GetChannelById(id int64) (*Channel, error) {
	if channel, ok := c.channelsById[id]; ok {
		return channel, nil
	}
	return nil, errors.New("No such channel")
}

/**
 * Delete a channel.
 */
func (c *ChannelService) DeleteChannel(channel *Channel) error {
	if _, ok := c.channelsById[channel.Id]; ok {
		delete(c.channelsById, channel.Id)
		return nil
	}
	return errors.New("No such channel")
}

/**
 * Lets a user to join a channel (if allowed)
 */
func (c *ChannelService) JoinChannel(channel *Channel, user *User) error {
	if channel.ContainsUser(user) {
		return nil
	}
	channel.Participants = append(channel.Participants, user)
	return c.SaveChannel(channel, true)
}

/**
 * Lets a user leave a channel or be kicked out.
 */
func (c *ChannelService) LeaveChannel(channel *Channel, user *User) error {
	for index, value := range channel.Participants {
		if value.Id == user.Id {
			channel.Participants = append(channel.Participants[:index], channel.Participants[index+1:]...)
			return c.SaveChannel(channel, true)
		}
	}
	return nil
}

/**
 * Returns the channels the user belongs to.
 */
func (c *ChannelService) ListChannels(user *User, team *Team) ([]*Channel, error) {
	out := make([]*Channel, 0, 100)
	for _, channel := range c.channelsById {
		if channel.Team.Id == team.Id {
			out = append(out, channel)
		}
	}
	return out, nil
}

/**
 * Removes all entries.
 */
func (svc *ChannelService) RemoveAllChannels() {
	svc.channelCounter = 1
	svc.channelsById = make(map[int64]*Channel)
}
