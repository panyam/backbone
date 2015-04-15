package gorilla

import (
	msgcore "github.com/panyam/relay/services/msg/core"
	// "log"
	"net/http"
)

// type ServiceRequestMaker func(*http.Request) (interface{}, error)
// type ServiceResponseMaker func(http.ResponseWriter, *http.Request, interface{}, error)

func GetTeamRequestMaker(request *http.Request) (interface{}, error) {
	req := msgcore.GetTeamRequest{}
	return &req, nil
}

func GetTeamResponsePresenter(http.ResponseWriter, *http.Request, interface{}, error) {
}

func SaveTeamRequestMaker(request *http.Request) (interface{}, error) {
	req := msgcore.SaveTeamRequest{}
	return &req, nil
}

func SaveTeamResponsePresenter(http.ResponseWriter, *http.Request, interface{}, error) {
}

func DeleteTeamRequestMaker(request *http.Request) (interface{}, error) {
	req := msgcore.DeleteTeamRequest{}
	return &req, nil
}

func DeleteTeamResponsePresenter(http.ResponseWriter, *http.Request, interface{}, error) {
}
