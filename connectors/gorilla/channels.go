package gorilla

import (
	"github.com/gorilla/mux"
	. "github.com/panyam/relay/connectors/gorilla/common"
	msgcore "github.com/panyam/relay/services/msg/core"
	"github.com/panyam/relay/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) GetChannelsHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		log.Println("Get Channels")
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

		teamId, _ := strconv.ParseInt(teamIdParam, 10, 64)
		team, _ := s.serviceGroup.TeamService.GetTeamById(teamId)
		participantsParam := strings.Split(request.FormValue("participants"), ",")
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
		s.serviceGroup.ChannelService.AddChannelMembers(channel, participantsParam)
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
