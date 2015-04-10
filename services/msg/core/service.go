package core

type ServiceGroup struct {
	IDService      IIDService
	ChannelService IChannelService
	TeamService    ITeamService
	UserService    IUserService
	MessageService IMessageService
}
