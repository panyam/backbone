package gorilla

import (
	"github.com/gorilla/mux"
	. "github.com/panyam/relay/connectors/gorilla/common"
	msgcore "github.com/panyam/relay/services/msg/core"
	"github.com/panyam/relay/utils"
	"log"
	"net/http"
	"strings"
)

func (s *Server) GetChannelsHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		teamIdParam := mux.Vars(request)["teamId"]
		teamId := utils.String2ID(teamIdParam)

		team, _ := s.serviceGroup.TeamService.GetTeamById(teamId)
		if team == nil {
			http.Error(rw, "No such team", http.StatusNotFound)
			return
		}

		ownerParam := request.FormValue("owner")
		var owner *msgcore.User = nil
		var err error = nil
		if ownerParam != "" {
			owner, err = s.serviceGroup.UserService.GetUser(ownerParam, team)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		matchall := request.FormValue("matchall") == "true"

		participantsParam := strings.Split(request.FormValue("participants"), ",")
		var participants []*msgcore.User = nil
		for _, participant := range participantsParam {
			user, err := s.serviceGroup.UserService.GetUser(participant, team)
			if err == nil {
				participants = append(participants, user)
			} else if matchall {
				// since this user does not exist then it cannot be matched
				// so return an empty list
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		order_by := request.FormValue("order_by")

		channels, members := s.serviceGroup.ChannelService.GetChannels(team, owner, order_by, participants, matchall)
		utils.SendJsonResponse(rw, map[string]interface{}{"channels": channels, "members": members})
	}
}

func (s *Server) CreateChannelHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		log.Println("Create Channels")
		teamIdParam := mux.Vars(request)["teamId"]
		if teamIdParam == "" {
			teamIdParam := request.FormValue("teamId")
			if teamIdParam == "" {
				http.Error(rw, "teamId not found", http.StatusBadRequest)
				return
			}
		}

		teamId := utils.String2ID(teamIdParam)
		team, _ := s.serviceGroup.TeamService.GetTeamById(teamId)
		if team == nil {
			http.Error(rw, "No such team", http.StatusNotFound)
			return
		}
		publicParam := request.FormValue("public")
		nameParam := request.FormValue("name")

		creator := context.Get("user").(*msgcore.User)
		channel := msgcore.NewChannel(team, creator, 0, nameParam)
		channel.Public = (publicParam != "false")
		err := s.serviceGroup.ChannelService.SaveChannel(channel, false)
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
	}
}

func (s *Server) ChannelDetailsHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		log.Println("GetChannelDetails")
	}
}

func (s *Server) UpdateChannelHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		log.Println("UpdateChannel")
	}
}

func (s *Server) DeleteChannelHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		log.Println("DeleteChannel")
	}
}
