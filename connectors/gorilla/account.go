package gorilla

import (
	"github.com/gorilla/mux"
	. "github.com/panyam/relay/connectors/gorilla/common"
	authcore "github.com/panyam/relay/services/auth/core"
	msgcore "github.com/panyam/relay/services/msg/core"
	"github.com/panyam/relay/utils"
	"log"
	"net/http"
	// "strconv"
)

func (s *Server) AccountRegisterHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		// teamId, _ := strconv.ParseInt(request.FormValue("teamId"), 10, 64)
		var team *msgcore.Team = nil
		// team, _ := s.serviceGroup.TeamService.GetTeamById(teamId)
		addressType := request.FormValue("address_type")
		verificationData := ""
		if addressType == "phone" {
			// generate a simple 5 digit PIN
			verificationData = utils.RandDigits(5)
		} else {
			verificationData = utils.RandAlnum(128)
		}
		log.Println("Verification Data: ", verificationData)
		registration := authcore.Registration{
			Username:         request.FormValue("username"),
			AddressType:      addressType,
			Address:          request.FormValue("address"),
			Team:             team,
			VerificationData: verificationData,
		}
		s.authService.SaveRegistration(&registration)
		utils.SendJsonResponse(rw, registration.ToDict())
	}
}

func (s *Server) AccountConfirmHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		registrationId := mux.Vars(request)["Id"]
		log.Println("Reg Id: ", registrationId)
	}
}

func (s *Server) AccountLoginHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		log.Println("Login...")
	}
}

func (s *Server) AccountLogoutHandler() RequestHandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request, context *RequestContext) {
		log.Println("Logout ...")
	}
}
