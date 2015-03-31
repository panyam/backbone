package middleware

import (
	"github.com/panyam/relay/models"
	"github.com/panyam/relay/services"
	"log"
	"net/http"
)

type UserService interface {
	/**
	 * Get user info by ID
	 */
	GetUserById(id string) (*User, error)
}

/**
 * Login based authenticator.
 * This checks that session ID and token values are valid and retrieves the user
 * based on this.
 */
type DebugValidator struct {
	userid      string
	userService UserService
}

func NewDebugValidator(u string, userService UserService) *DebugValidator {
	return &DebugValidator{userid: u}
}

func (validator *DebugValidator) ValidateRequest(request *http.Request, context *RequestContext) *models.User {
	if request.Form["__user"] == validator.Userid {
		return userService.GetUserById(validator.Userid)
	}
	return nil
}
