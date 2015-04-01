package goclient

import (
	"encoding/json"
	"fmt"
	authcore "github.com/panyam/relay/services/auth/core"
	msgcore "github.com/panyam/relay/services/msg/core"
	"io"
	"log"
	"net/http"
	"strings"
	// "log"
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

func (client *ApiClient) RegisterUser(teamId int64, username string, address_type string, address string, password string) (*authcore.Registration, error) {
	params := map[string]string{
		"teamId":       fmt.Sprintf("%d", teamId),
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

	log.Println("data: ", data)
	registration := authcore.Registration{}
	registration.Username = data["Username"].(string)
	registration.Address = data["Address"].(string)
	registration.AddressType = data["AddressType"].(string)
	registration.Id, _ = data["Id"].(json.Number).Int64()
	return &registration, err
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

func (client *ApiClient) ConfirmRegistration(registrationId string, verificationCode string) error {
	params := map[string]string{"verification_code": verificationCode}
	url := client.url(fmt.Sprintf("/registrations/%s/confirm/", registrationId))
	req, err := MakeRequest("POST", url, params, nil)
	resp, err := SendRequest(req, nil)
	log.Println("Confirm Reg Resp: ", resp)
	defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	return err
}

func (client *ApiClient) GetTeams(username string) ([]*msgcore.Team, error) {
	url := fmt.Sprintf("/users/%s/teams/", username)
	req, err := client.MakeAuthRequest("GET", url, nil, nil)
	teams := new([]*msgcore.Team)
	_, err = SendRequest(req, teams)
	return *teams, err
}

func (client *ApiClient) CreateTeam(name string, organization string) (*msgcore.Team, error) {
	params := map[string]string{"name": name, "org": organization}
	req, err := client.MakeAuthRequest("POST", "/teams/", params, nil)
	team := new(msgcore.Team)
	_, err = SendRequest(req, team)
	return team, err
}
