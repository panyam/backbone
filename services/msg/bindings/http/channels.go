package http

import (
	msgcore "github.com/panyam/relay/services/msg/core"
	"net/http"
)

HTTPBindingFor(Channel) - {
		/teams/{teamId}/
		"id" -> Channel.Id
}

/**
 * Channels/Threads/Groups/Conversations etc
 */
func ChannelToHttpRequest(c *msgcore.Channel) (*http.Request, error) {
	return nil, nil
}
