package middleware

import (
	"github.com/panyam/relay/connectors/gorilla/common"
	msgcore "github.com/panyam/relay/services/msg/core"
	"github.com/panyam/relay/utils"
	"net/http"
)

type UserService interface {
	/**
	 * Get user info by ID
	 */
	GetUserById(id int64) (*msgcore.User, error)
}

/**
 * Login based authenticator.
 * This checks that session ID and token values are valid and retrieves the user
 * based on this.
 */
type DebugValidator struct {
	userid      int64
	userService UserService
}

func NewDebugValidator(u int64, userSvc UserService) *DebugValidator {
	return &DebugValidator{userid: u, userService: userSvc}
}

func (validator *DebugValidator) ValidateRequest(request *http.Request, context common.IRequestContext) *msgcore.User {
	users := request.Form["__dbguser"]
	if len(users) > 0 {
		userid := utils.String2ID(users[0])
		if userid == validator.userid {
			user, err := validator.userService.GetUserById(validator.userid)
			if user != nil && err == nil {
				return user
			}
		}
	}
	return nil
}
