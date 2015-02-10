package services

import (
	"errors"
	"fmt"
	. "github.com/panyam/backbone/models"
)

type ChannelService struct {
	Cls            IChannelService
	channelsByName map[string]*Channel
}

func NewChannelService() *ChannelService {
	svc := ChannelService{}
	svc.Cls = &svc
	svc.channelsByName = make(map[string]*Channel)
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
 * Retrieve a channel by ID.
 */
func (c *ChannelService) GetChannelById(channelId string) (*Channel, error) {
	return nil, nil
}

/**
 * Delete a channel.
 */
func (c *ChannelService) DeleteChannel(channel *Channel) error {
	return nil
}

/**
 * Returns the channels the user belongs to.
 */
func (c *ChannelService) ListChannels(user *User) ([]*Channel, error) {
	return nil, nil
}

/**
 * Lets a user to join a channel (if allowed)
 */
func (c *ChannelService) JoinChannel(user *User, channel *Channel) error {
	return nil
}

/**
 * Lets a user leave a channel or be kicked out.
 */
func (c *ChannelService) LeaveChannel(user *User, channel *Channel, forced bool) error {
	return nil
}

/**
 * Invite a user to a channel.
 */
func (c *ChannelService) InviteToChannel(inviter *User, invitee *User, channel *Channel) error {
	return nil
}
