package service

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
	 */
	CreateUser(id string, username string) (*User, error)
}

type ITeamService interface {
	/**
	 * Retrieve a team by ID.
	 */
	GetTeamById(teamId string) (Team, error)

	/**
	 * Lets a user create a team.
	 */
	CreateTeam(teamName string) (Team, error)

	/**
	 * Delete a team.
	 */
	DeleteTeam(team Team) error

	/**
	 * Returns the teams the user belongs to.
	 */
	ListTeams(user User) ([]Team, error)

	/**
	 * Lets a user to join a team (if allowed)
	 */
	JoinTeam(user User, team Team) error

	/**
	 * Lets a user leave a team or be kicked out.
	 */
	LeaveTeam(user User, team Team, forced bool) error

	/**
	 * Invite a user to a team.
	 */
	InviteToTeam(user User, team Team) error
}

type IChannelService interface {
	/**
	 * Lets a user create a channel.
	 */
	CreateChannel(channelName string) (Channel, error)

	/**
	 * Retrieve a channel by ID.
	 */
	GetChannelById(channelId string) (Channel, error)

	/**
	 * Delete a channel.
	 */
	DeleteChannel(channel Channel) error

	/**
	 * Returns the channels the user belongs to.
	 */
	ListChannels(user User) ([]Channel, error)

	/**
	 * Lets a user to join a channel (if allowed)
	 */
	JoinChannel(user User, channel Channel) error

	/**
	 * Lets a user leave a channel or be kicked out.
	 */
	LeaveChannel(user User, channel Channel, forced bool) error

	/**
	 * Invite a user to a channel.
	 */
	InviteToChannel(inviter User, invitee User, channel Channel) error
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
