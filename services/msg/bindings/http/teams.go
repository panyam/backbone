package http

import (
	msgcore "github.com/panyam/relay/services/msg/core"
	"net/http"
)

func MakeRequest(method, endpoint string,
	queryParams map[string]string,
	body io.Reader) (*http.Request, error) {
	if queryParams != nil {
		i := 0
		for key, value := range queryParams {
			if i == 0 {
				endpoint += "?"
			} else {
				endpoint += "&"
			}
			endpoint += key
			endpoint += "=" + value
			i++
		}
	}
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		log.Printf("%s %s Error: %s\n", method, endpoint, err)
		return nil, err
	}
	return req, nil
}

func SaveTeamRequestToHttp(input *Team) (*http.Request, error) {
	url := fmt.Sprintf("/users/teams/?Offset=%d&Count=%d&Organization=%s",
		request.Offset, request.Count, request.Organization)
	return client.MakeRequest("GET", url, nil, nil), nil
}

func GetTeamsRequestToHttp(request *GetTeamsRequest) (*http.Request, error) {
}

func GetTeamRequestToHttp(team *Team) (*http.Request, error) {
}

func DeleteTeamRequestToHttp(team *Team) (*http.Request, error) {
}

func JoinTeamRequestToHttp(request *TeamMembershipRequest) (*http.Request, error) {
}

func TeamContainsRequestToHttp(request *TeamMembershipRequest) (*http.Request, error) {
}

func LeaveTeamRequestToHttp(request *TeamMembershipRequest) (*http.Request, error) {
}

func SaveTeamRequestFromHttp(*http.Request) (*Team, error) {
	url := fmt.Sprintf("/users/teams/?Offset=%d&Count=%d&Organization=%s",
		request.Offset, request.Count, request.Organization)
	return client.MakeRequest("GET", url, nil, nil), nil
}

func GetTeamsRequestFromHttp(request *http.Request) (*GetTeamsRequest, error) {
}

func GetTeamRequestFromHttp(request *http.Request) (*Team, error) {
}

func DeleteTeamRequestFromHttp(request *http.Request) (*Team, error) {
}

func JoinTeamRequestFromHttp(request *http.Request) (*TeamMembershipRequest, error) {
}

func TeamContainsRequestFromHttp(request *http.Request) (*TeamMembershipRequest, error) {
}

func LeaveTeamRequestFromHttp(request *http.Request) (*TeamMembershipRequest, error) {
}
