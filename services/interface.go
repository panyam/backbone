package services

import (
	. "github.com/panyam/backbone/models"
)

/**
 * Base service operations.  These dont care about authorization for now and
 * assume the user is authorized.  Authn (and possible Authz) have to be taken
 * care of seperately.
 */
type IUserService interface {
	/**
	 * Get user info by ID
	 */
	GetUserById(id string) (*User, error)

	/**
	 * Get a user by username.
	 */
	GetUser(username string) (*User, error)

	/**
	 * Create a user with the given id and username.
	 * If the ID or Username already exists an error is thrown.
	 */
	CreateUser(id string, username string) (*User, error)

	/**
	 * Saves a user details.
	 * If the user id or username does not exist an error is thrown.
	 * If the username or user id already exist and are not the same
	 * object then an error is thrown.
	 */
	SaveUser(user *User) error
}

type IChannelService interface {
	/**
	 * Lets a user create a channel.  If a channel exists an error is returned.
	 */
	CreateChannel(channelName string) (*Channel, error)

	/**
	 * Retrieve a channel by name.
	 */
	GetChannelByName(name string) (*Channel, error)

	/**
	 * Delete a channel.
	 */
	DeleteChannel(channel *Channel) error

	/**
	 * Returns the channels the user belongs to.
	 */
	ListChannels(user *User) ([]*Channel, error)

	/**
	 * Lets a user to join a channel (if allowed)
	 */
	JoinChannel(channel *Channel, user *User) error

	/**
	 * Tells if a user belongs to a channel.
	 */
	ChannelContains(channel *Channel, user *User) bool

	/**
	 * Lets a user leave a channel or be kicked out.
	 */
	LeaveChannel(channel *Channel, user *User) error

	/**
	 * Invite a user to a channel.
	 */
	InviteToChannel(inviter *User, invitee *User, channel *Channel) error
}

type IMessageService interface {
	/**
	 * Get the messages in a channel for a particular user.
	 */
	GetMessages(channel Channel, user User, offset int, count int) ([]Message, error)

	/**
	 * Send messages to particular recipients in this channel.
	 */
	SendMessage(channel Channel, message Message, sender User, recipients []User) error

	/**
	 * Remove a particular message.
	 */
	DeleteMessage(message Message) error
}
