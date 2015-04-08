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

		channel := msgcore.NewChannel(team, context.Get("user").(*msgcore.User), 0, nameParam)
		channel.Public = (publicParam != "false")
		s.serviceGroup.ChannelService.SaveChannel(channel, false)

		participantIDs := make([]int64, 0, len(participantsParam))
		for index, participantID := range participantsParam {
			log.Println("Index, partID: ", index, participantID, participantIDs)
			partID, err := strconv.ParseInt(participantID, 10, 64)
			if err == nil {
				participantIDs = append(participantIDs, partID)
			}
		}
		s.serviceGroup.ChannelService.AddChannelMembers(channel, participantIDs)
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
