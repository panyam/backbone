package services

import (
	"errors"
	"fmt"
	. "github.com/panyam/backbone/models"
)

type ChannelService struct {
	Cls            IChannelService
	channelsByName map[string]*Channel
	usersById      map[string]*User
}

func NewChannelService() *ChannelService {
	svc := ChannelService{}
	svc.Cls = &svc
	svc.channelsByName = make(map[string]*Channel)
	svc.usersById = make(map[string]*User)
	return &svc
}

/**
 * Lets a user create a channel.
 */
func (c *ChannelService) CreateChannel(channelName string) (*Channel, error) {
	if _, ok := c.channelsByName[channelName]; ok {
		return nil, errors.New("Channel already exists")
	}
	channel := Channel{Name: channelName}
	channel.Id = fmt.Sprintf("%d", len(c.channelsByName))
	c.channelsByName[channelName] = &channel
	return &channel, nil
}

/**
 * Retrieve a channel by Name.
 */
func (c *ChannelService) GetChannelByName(channelName string) (*Channel, error) {
	if channel, ok := c.channelsByName[channelName]; ok {
		return channel, nil
	}
	return nil, errors.New("No such channel")
}

/**
 * Delete a channel.
 */
func (c *ChannelService) DeleteChannel(channel *Channel) error {
	if channel, ok := c.channelsByName[channel.Name]; ok {
		delete(c.channelsByName, channel.Name)
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
