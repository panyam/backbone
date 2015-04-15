package gorilla

import (
	// "github.com/gorilla/mux"
	msgcore "github.com/panyam/relay/services/msg/core"
	// "github.com/panyam/relay/utils"
	// "log"
	"net/http"
	// "strings"
)

func GetChannelRequestMaker(request *http.Request) (interface{}, error) {
	req := msgcore.Channel{}
	return &req, nil
}

func GetChannelResponsePresenter(rw http.ResponseWriter, req *http.Request, result interface{}, err error) {
	DefaultResponsePresenter(rw, req, result, err)
}

func GetChannelsRequestMaker(request *http.Request) (interface{}, error) {
	req := msgcore.GetChannelsRequest{}
	return &req, nil
}

func GetChannelsResponsePresenter(rw http.ResponseWriter, req *http.Request, result interface{}, err error) {
	DefaultResponsePresenter(rw, req, result, err)
}

func CreateChannelRequestMaker(request *http.Request) (interface{}, error) {
	req := msgcore.CreateChannelRequest{}
	return &req, nil
	/*
		teamIdParam := mux.Vars(request)["teamId"]
		if teamIdParam == "" {
			teamIdParam := request.FormValue("teamId")
			if teamIdParam == "" {
				http.Error(rw, "teamId not found", http.StatusBadRequest)
				return
			}
		}

		teamId := utils.String2ID(teamIdParam)
		team, _ := s.serviceGroup.TeamService.GetTeam(msgcore.NewTeamById(teamId))
		if team == nil {
			http.Error(rw, "No such team", http.StatusNotFound)
			return
		}
		publicParam := request.FormValue("public")
		nameParam := request.FormValue("name")

		creator := context.Get("user").(*msgcore.User)
		channel := msgcore.NewChannel(team, creator, 0, nameParam)
		channel.Public = (publicParam != "false")
		err := s.serviceGroup.ChannelService.CreateChannel(channel, false)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		s.serviceGroup.ChannelService.AddChannelMembers(channel, []string{creator.Username})

		participantsParam := strings.Split(request.FormValue("participants"), ",")
		if len(participantsParam) > 0 {
			s.serviceGroup.ChannelService.AddChannelMembers(channel, participantsParam)
		}
		utils.SendJsonResponse(rw, channel.ToDict())
	*/
}

func CreateChannelResponsePresenter(rw http.ResponseWriter, req *http.Request, result interface{}, err error) {
	DefaultResponsePresenter(rw, req, result, err)
}

func UpdateChannelRequestMaker(request *http.Request) (interface{}, error) {
	req := msgcore.Channel{}
	return &req, nil
}

func UpdateChannelResponsePresenter(rw http.ResponseWriter, req *http.Request, result interface{}, err error) {
	DefaultResponsePresenter(rw, req, result, err)
}

func DeleteChannelRequestMaker(request *http.Request) (interface{}, error) {
	req := msgcore.DeleteChannelRequest{}
	return &req, nil
}

func DeleteChannelResponsePresenter(rw http.ResponseWriter, req *http.Request, result interface{}, err error) {
	DefaultResponsePresenter(rw, req, result, err)
}
