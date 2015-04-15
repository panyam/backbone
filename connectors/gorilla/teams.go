package gorilla

import (
	"errors"
	"github.com/gorilla/mux"
	msgcore "github.com/panyam/relay/services/msg/core"
	"github.com/panyam/relay/utils"
	// "log"
	"net/http"
)

func DefaultResponsePresenter(rw http.ResponseWriter, req *http.Request, result interface{}, err error) {
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	} else {
		utils.SendJsonResponse(rw, result)
	}
}

// type ServiceRequestMaker func(*http.Request) (interface{}, error)
// type ServiceResponseMaker func(http.ResponseWriter, *http.Request, interface{}, error)

func GetTeamRequestMaker(request *http.Request) (interface{}, error) {
	req := msgcore.Team{}

	teamIdParam := mux.Vars(request)["teamId"]
	if teamIdParam == "" {
		teamIdParam := request.FormValue("teamId")
		if teamIdParam == "" {
			return nil, errors.New("Team ID Not found")
		}
	}

	req.Id = utils.String2ID(teamIdParam)
	return &req, nil
}

func GetTeamResponsePresenter(rw http.ResponseWriter, req *http.Request, result interface{}, err error) {
	DefaultResponsePresenter(rw, req, result, err)
}

func SaveTeamRequestMaker(request *http.Request) (interface{}, error) {
	req := msgcore.Team{}
	return &req, nil
}

func SaveTeamResponsePresenter(rw http.ResponseWriter, req *http.Request, result interface{}, err error) {
	DefaultResponsePresenter(rw, req, result, err)
}

func DeleteTeamRequestMaker(request *http.Request) (interface{}, error) {
	req := msgcore.Team{}
	return &req, nil
}

func DeleteTeamResponsePresenter(rw http.ResponseWriter, req *http.Request, result interface{}, err error) {
	DefaultResponsePresenter(rw, req, result, err)
}
