package goclient

import (
	"errors"
	"fmt"
	authcore "github.com/panyam/relay/services/auth/core"
	msgcore "github.com/panyam/relay/services/msg/core"
	"github.com/panyam/relay/utils"
	"io"
	"log"
	"net/http"
	"strings"
	// "os"
	// "os/exec"
	// "sort"
	// "strconv"
)

/**
 * Creates a request that can be sent to register a new user.
 */
func (client *ApiClient) NewRegisterUserRequest(team *msgcore.Team, username string, address_type string, address string, password string) (*http.Request, error) {
	params := map[string]string{
		"teamId":       utils.ID2String(team.Id),
		"username":     username,
		"address_type": address_type,
		"address":      address,
	}
	if password != "" {
		params["password"] = password
	}
	req, err := MakeRequest("POST", client.url("/users/register/"), params, nil)
	return req, err
}

/**
 * Creates a request to call the registration confirmation API.
 */
func (client *ApiClient) NewConfirmRegistrationRequest() (*http.Request, error) {
	params := map[string]string{"verification_code": verificationCode}
	url := client.url(fmt.Sprintf("/users/registrations/%s/confirm/", utils.ID2String(registration.Id)))
	return MakeRequest("POST", url, params, nil)
}

// Gets a list of teams a user is subscribed or invited to.
//
// **Endpoints:** GET /users/teams/
//
// **Auth Required:** YES
//
// **Return:**
//  A list of teams that the current user is subscribed to.  eg:
//  	```
//  	[ {"id": 123, "name": "Dream Team", "organization": "Dream Owner"} ]
//  	```
func (client *ApiClient) NewGetTeamsRequest() (*http.Request, error) {
	url := "/users/teams/"
	return client.MakeAuthRequest("GET", url, nil, nil)
}

// Create a team
//
// **Endpoints:** POST /teams/
//
// **Auth Required:** YES
//
// **Parameters:**
// - organization: Organization the team belongs to (optional)
// - name: Name of the team (required and must be unique within the organization).
//
// **Return:**
//
// HTTP Status 200 on success along with team details, eg:
//
// ```
// {"id": "123", "name": "Dream Team", "organization": "Dream Owner"}
// ```
//
// The user needs to have the "teamcreator" permission set otherwise returns
// 4xx.
func (client *ApiClient) NewCreateTeamRequest(team *Team) (*http.Request, error) {
	params := map[string]string{"name": team.Name, "org": team.Organization}
	return client.MakeAuthRequest("POST", "/teams/", params, nil)
}

// Create a channel
//
// **Endpoints:** POST /teams/<id>/channels/
//
// **Auth Required:** YES and user must belong to team or have
// createteam_<teamid> permission.
//
// **Parameters:**
//
// - name: Name of the channel (required)
// - participants: Comma seperated list of usernames
// - public: Whether the channel is public or not (public by default)
//
// **Return:**
//
// HTTP Status 200 on success along with team details, eg:
//
// ```
// {"id": "123", "name": "Dream Team", "organization": "Dream Owner"}
// ```
func (client *ApiClient) CreateChannel(request *CreateChannelRequest) (*msgcore.Channel, error) {
	params := map[string]interface{}{
		"name":         request.Channel.Name,
		"public":       true,
		"participants": participants,
	}
	if !public {
		params["public"] = false
	}
	url := fmt.Sprintf("/teams/%s/channels/", utils.ID2String(request.Channel.Team.Id))
	return client.MakeAuthRequest("POST", url, params, nil)
}
