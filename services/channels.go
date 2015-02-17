package services

import (
	"errors"
	"fmt"
	. "github.com/panyam/backbone/models"
)

type ChannelService struct {
	Cls            IChannelService
	channelCounter int64
	channelsByKey  map[string]*Channel
	usersById      map[string]*User
}

func NewChannelService() *ChannelService {
	svc := ChannelService{}
	svc.Cls = &svc
	svc.channelCounter = 1
	svc.channelsByKey = make(map[string]*Channel)
	svc.usersById = make(map[string]*User)
	return &svc
}

/**
 * Lets a user create a channel.
 */
func (c *ChannelService) CreateChannel(id string, channelGroup string, channelName string) (*Channel, error) {
	key := channelGroup + ":" + channelName
	if _, ok := c.channelsByKey[key]; ok {
		return nil, errors.New("Channel already exists")
	}
	if id == "" {
		id = fmt.Sprintf("%d", c.channelCounter)
	}
	channel := NewChannel(id, channelGroup, channelName)
	c.channelsByKey[key] = channel
	c.channelCounter++
	return channel, nil
}

/**
 * Retrieve channels in a group
 */
func (c *ChannelService) GetChannelsInGroup(group string, offset int, count int) ([]*Channel, error) {
	return nil, nil
}

/**
 * Retrieve a channel by Name.
 */
func (c *ChannelService) GetChannelByName(channelGroup string, channelName string) (*Channel, error) {
	key := channelGroup + ":" + channelName
	if channel, ok := c.channelsByKey[key]; ok {
		return channel, nil
	}
	return nil, errors.New("No such channel")
}

/**
 * Delete a channel.
 */
func (c *ChannelService) DeleteChannel(channel *Channel) error {
	key := channel.Group + ":" + channel.Name
	if _, ok := c.channelsByKey[key]; ok {
		delete(c.channelsByKey, key)
		return nil
	}
	return errors.New("No such channel")
}

/**
 * Lets a user to join a channel (if allowed)
 */
func (c *ChannelService) JoinChannel(channel *Channel, user *User) error {
	c.usersById[user.Id] = user
	return nil
}

/**
 * Tells if a user belongs to a channel.
 */
func (c *ChannelService) ChannelContains(channel *Channel, user *User) bool {
	_, ok := c.usersById[user.Id]
	return ok
}

/**
 * Returns the channels the user belongs to.
 */
func (c *ChannelService) ListChannels(user *User) ([]*Channel, error) {
	return nil, nil
}

/**
 * Lets a user leave a channel or be kicked out.
 */
func (c *ChannelService) LeaveChannel(channel *Channel, user *User) error {
	delete(c.usersById, user.Id)
	return nil
}

/**
 * Invite a user to a channel.
 */
func (c *ChannelService) InviteToChannel(inviter *User, invitee *User, channel *Channel) error {
	return nil
}
