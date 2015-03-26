package sqlds

import (
	"database/sql"
	"errors"
	. "github.com/panyam/backbone/services/core"
)

type ChannelService struct {
	Cls interface{}
	DB  *sql.DB
}

func NewChannelService(db *sql.DB) *ChannelService {
	svc := ChannelService{}
	svc.Cls = &svc
	svc.DB = db
	svc.RemoveAllChannels()
	return &svc
}

/**
 * Lets a user create a channel.
 */
func (c *ChannelService) SaveChannel(channel *Channel, override bool) error {
	return nil
}

/**
 * Retrieve a channel by Name.
 */
func (c *ChannelService) GetChannelById(id int64) (*Channel, error) {
	return nil, errors.New("No such channel")
}

/**
 * Delete a channel.
 */
func (c *ChannelService) DeleteChannel(channel *Channel) error {
	return errors.New("No such channel")
}

/**
 * Lets a user to join a channel (if allowed)
 */
func (c *ChannelService) JoinChannel(channel *Channel, user *User) error {
	return nil
}

/**
 * Lets a user leave a channel or be kicked out.
 */
func (c *ChannelService) LeaveChannel(channel *Channel, user *User) error {
	return nil
}

/**
 * Returns the channels the user belongs to.
 */
func (c *ChannelService) ListChannels(user *User, team *Team) ([]*Channel, error) {
	return nil, nil
}

/**
 * Removes all entries.
 */
func (svc *ChannelService) RemoveAllChannels() {
}
