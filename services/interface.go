package services

import (
	. "github.com/panyam/backbone/models"
)

type IIDService interface {
	/**
	 * Creates a new ID.
	 */
	CreateID(domain string) string

	/**
	 * Releases an ID back to the domain.
	 */
	ReleaseID(domain string, id string)
}

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
	 * Saves a user details.
	 * If the user id or username does not exist an error is thrown.
	 * If the username or user id already exist and are not the same
	 * object then an error is thrown.
	 */
	SaveUser(user *User) error

	/**
	 * Deletes a user from the sytem
	 */
	// DeleteUser(user *User) error

	/**
	 * Create a user with the given id and username.
	 * If the ID or Username already exists an error is thrown.
	 * If the ID is empty, then it is upto the backend to decide whether to
	 * throw an error or auto assign an ID.
	 * A valid User object on return WILL have an ID if the backend can
	 * auto generate IDs
	 */
	CreateUser(id string, username string) (*User, error)
}

type ITeamService interface {
	/**
	 * Create a team.
	 * If the ID is empty, then it is upto the backend to decide whether to
	 * throw an error or auto assign an ID.
	 * A valid Team object on return WILL have an ID if the backend can
	 * auto generate IDs
	 */
	CreateTeam(id string, org string, name string) (*Team, error)

	/**
	 * Retrieve teams in a org
	 */
	GetTeamsInOrg(org string, offset int, count int) ([]*Team, error)

	/**
	 * Retrieve a team by Name.
	 */
	GetTeamByName(org string, name string) (*Team, error)

	/**
	 * Delete a team.
	 */
	DeleteTeam(team *Team) error

	/**
	 * Lets a user to join a team (if allowed)
	 */
	JoinTeam(team *Team, user *User) error

	/**
	 * Tells if a user belongs to a team.
	 */
	TeamContains(team *Team, user *User) bool

	/**
	 * Lets a user leave a team or be kicked out.
	 */
	LeaveTeam(team *Team, user *User) error
}

type IChannelService interface {
	/**
	 * Lets a user create a channel.  If a channel exists an error is returned.
	 * If the ID is empty, then it is upto the backend to decide whether to
	 * throw an error or auto assign an ID.
	 * A valid Channel object on return WILL have an ID if the backend can
	 * auto generate IDs
	 */
	CreateChannel(id string, channelGroup string, channelName string) (*Channel, error)

	/**
	 * Retrieve all channels in a particular group.
	 */
	GetChannelsInGroup(group string, offset int, count int) ([]*Channel, error)

	/**
	 * Retrieve a channel by name within a particular group if it exists.
	 */
	GetChannelByName(group string, name string) (*Channel, error)

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
	GetMessages(channel *Channel, user *User, offset int, count int) ([]*Message, error)

	/**
	 * Creates a message to particular recipients in this channel.  This is
	 * called "Create" instead of "Send" so as to not confuse with the delivery
	 * details.
	 * If message ID is empty then the backend can auto generate one if it is
	 * capable of doing so.
	 * A valid Message object on return WILL have a non empty ID if the backend can
	 * auto generate IDs
	 */
	CreateMessage(message *Message) error

	/**
	 * Remove a particular message.
	 */
	DeleteMessage(message *Message) error

	/**
	 * Saves a message.
	 * If the message ID is missing (or empty) then a new message is created.
	 * If message ID is present then the existing message is updated.
	 */
	// SaveMessage(message *Message) error
}
