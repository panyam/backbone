package connectors

type HTTPConnector struct {
	userService    IUserService
	channelService IChannelService
	messageService IMessageService
}
