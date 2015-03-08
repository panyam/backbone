package core

type ServiceGroup struct {
	ChannelService IChannelService
	TeamService    ITeamService
	UserService    IUserService
	MessageService IMessageService
}
