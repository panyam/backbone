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

type ApiClient struct {
	Url           string
	Authenticator Authenticator
}

func NewApiClient(url string) *ApiClient {
	client := ApiClient{Url: url}
	log.Println("New Api Client, URL: ", url)
	return &client
}

func (client *ApiClient) url(path string) string {
	out := client.Url
	if strings.HasPrefix(path, "/") {
		return out + path
	} else {
		return out + "/" + path
	}
}

func (client *ApiClient) EnableAuthentication(auth Authenticator) {
	client.Authenticator = auth
}

func (client *ApiClient) DisableAuthentication() {
	client.Authenticator = nil
}

func (client *ApiClient) MakeAuthRequest(method, endpoint string,
	queryParams map[string]string, body io.Reader) (*http.Request, error) {
	request, err := MakeRequest(method, client.url(endpoint), queryParams, body)
	if err != nil {
		return nil, err
	}
	if client.Authenticator != nil {
		client.Authenticator.AuthenticateRequest(request)
	}
	return request, nil
}

func (client *ApiClient) RegisterUser(team *msgcore.Team, username string, address_type string, address string, password string) (*authcore.Registration, error) {
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
	var data map[string]interface{}
	_, err = SendRequest(req, &data)
	if err != nil {
		return nil, err
	}

	return authcore.RegistrationFromDict(data)
}

func (client *ApiClient) ConfirmRegistration(registration *authcore.Registration, verificationCode string) error {
	params := map[string]string{"verification_code": verificationCode}
	url := client.url(fmt.Sprintf("/users/registrations/%s/confirm/", utils.ID2String(registration.Id)))
	req, err := MakeRequest("POST", url, params, nil)
	resp, err := SendRequest(req, nil)
	if err != nil {
		if resp.StatusCode == 400 {
			return errors.New("Confirmation failed")
		} else if resp.StatusCode != 200 {
			return errors.New(resp.Status)
		}
	}
	return err
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
func (client *ApiClient) GetTeams() ([]*msgcore.Team, error) {
	url := "/users/teams/"
	req, err := client.MakeAuthRequest("GET", url, nil, nil)
	var data map[string]interface{}
	_, err = SendRequest(req, data)
	return nil, err
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
func (client *ApiClient) CreateTeam(name string, organization string) (*msgcore.Team, error) {
	params := map[string]string{"name": name, "org": organization}
	req, err := client.MakeAuthRequest("POST", "/teams/", params, nil)
	team := new(msgcore.Team)
	resp, err := SendRequest(req, team)
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	return team, err
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
// - participants: List of usernames participating in the channel.
// - public: Whether the channel is public or not (public by default)
//
// **Return:**
//
// HTTP Status 200 on success along with team details, eg:
//
// ```
// {"id": "123", "name": "Dream Team", "organization": "Dream Owner"}
// ```
func (client *ApiClient) CreateChannel(team *msgcore.Team, name string, public bool, participants []string) (*msgcore.Channel, error) {
	params := map[string]string{"name": name,
		"participants": strings.Join(participants, ","),
		"public":       "true",
	}
	if !public {
		params["public"] = "false"
	}
	url := fmt.Sprintf("/teams/%s/channels/", utils.ID2String(team.Id))
	req, err := client.MakeAuthRequest("POST", url, params, nil)
	var data map[string]interface{}
	resp, err := SendRequest(req, data)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	return msgcore.ChannelFromDict(data)
}
