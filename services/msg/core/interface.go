package core

import ()

type IIDService interface {
	/**
	 * Creates a new ID.
	 */
	CreateID(domain string) string

	/**
	 * Releases an ID back to the domain.
	 */
	ReleaseID(domain string, id int64)
}

/**
 * Base service operations.  These dont care about authorization for now and
 * assume the user is authorized.  Authn (and possible Authz) have to be taken
 * care of seperately.
 */
type IUserService interface {
	/**
	 * Removes all entries.
	 */
	RemoveAllUsers()

	/**
	 * Get user info by ID
	 */
	GetUserById(id int64) (*User, error)

	/**
	 * Get a user by username in a particular team.
	 */
	GetUser(username string, team *Team) (*User, error)

	/**
	 * Saves a user.
	 * 	If the ID param is empty:
	 * 		If username/team does not already exist a new one is created.
	 * 		otherwise, it is updated and returned if override=true otherwise
	 * 		false is returned.
	 * 	Otherwise:
	 * 		If username/team does not exist then it is written as is (Create or Update)
	 * 		otherwise if IDs of curr and existing are different errow is thrown,
	 * 		otherwise object is updated.
	 */
	SaveUser(user *User, override bool) error
}

type ITeamService interface {
	/**
	 * Removes all entries.
	 */
	RemoveAllTeams()

	/**
	 * Create a team.
	 * If the ID is empty, then it is upto the backend to decide whether to
	 * throw an error or auto assign an ID.
	 * A valid Team object on return WILL have an ID if the backend can
	 * auto generate IDs
	 */
	CreateTeam(id int64, org string, name string) (*Team, error)

	/**
	 * Retrieve teams in a org
	 */
	GetTeamsInOrg(org string, offset int, count int) ([]*Team, error)

	/**
	 * Retrieve a team by Id
	 */
	GetTeamById(id int64) (*Team, error)

	/**
	 * Retrieve a team by Name.
	 */
	GetTeamByName(org string, name string) (*Team, error)

	/**
	 * Delete a team.
	 */
	DeleteTeam(team *Team) error

	/**
	 * Lets a user with the given username join a team (if allowed)
	 */
	JoinTeam(team *Team, username string) (*User, error)

	/**
	 * Tells if a user belongs to a team.
	 */
	TeamContains(team *Team, username string) bool

	/**
	 * Lets a user leave a team or be kicked out.
	 */
	LeaveTeam(team *Team, user *User) error
}

type IChannelService interface {
	/**
	 * Removes all entries.
	 */
	RemoveAllChannels()

	/**
	 * Creates a channel.
	 * If the channel's ID parameter is not set then a new channel is created.
	 * If the ID parameter IS set:
	 * 		if override parameter is true, the channel is upserted (updated if it
	 * 		existed, otherwise created).
	 * 		If the override parameter is false, then the channel is only saved
	 * 		if it does not already exist and returns a ChannelExists error if an
	 * 		existing channel with the same ID exists.
	 */
	SaveChannel(channel *Channel, override bool) error

	/**
	 * Get channel by Id
	 */
	GetChannelById(id int64) (*Channel, error)

	/**
	 * Gets the channel members.
	 */
	GetChannelMembers(channel *Channel) []ChannelMember

	/**
	 * Tells if a channe contains a user.
	 */
	ContainsUser(channel *Channel, user *User) bool

	/**
	 * Delete a channel.
	 */
	DeleteChannel(channel *Channel) error

	/**
	 * Returns the channels the user belongs to in a given team.
	 */
	ListChannels(user *User, team *Team) ([]*Channel, error)

	/**
	 * Lets a user to join a channel (if allowed)
	 */
	JoinChannel(channel *Channel, user *User) error

	/**
	 * Lets a user leave a channel or be kicked out.
	 */
	LeaveChannel(channel *Channel, user *User) error
}

type IMessageService interface {
	/**
	 * Removes all entries.
	 */
	RemoveAllMessages()

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
	SaveMessage(message *Message) error

	/**
	 * Gets a message by ID
	 */
	GetMessageById(id int64) (*Message, error)

	/**
	 * Gets the fragments of a message.
	 */
	GetMessageFragments(message *Message) []*MessageFragment

	/**
	 * Get receipts of a particular message.
	 */
	GetMessageReceipts(message *Message) []*MessageReceipt

	/**
	 * Remove a particular message.
	 */
	DeleteMessage(message *Message) error
}
