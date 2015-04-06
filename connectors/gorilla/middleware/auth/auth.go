package middleware

import (
	"errors"
	"github.com/panyam/relay/connectors/gorilla/common"
	msgcore "github.com/panyam/relay/services/msg/core"
	"log"
	"net/http"
)

type AuthValidator interface {
	ValidateRequest(request *http.Request, context common.IRequestContext) *msgcore.User
}

type AuthMiddleware struct {
	Validators []AuthValidator
}

func (am *AuthMiddleware) ProcessRequest(rw http.ResponseWriter, request *http.Request, context common.IRequestContext) error {
	for _, validator := range am.Validators {
		user := validator.ValidateRequest(request, context)
		if user != nil {
			context.Set("user", user)
			return nil
		}
	}

	// none of the validators matched so throw a 401 unauthorized response
	log.Println("Auth Failed")
	rw.WriteHeader(http.StatusUnauthorized)
	return errors.New("Not allowed")
}

func (validator *AuthMiddleware) ProcessResponse(rw http.ResponseWriter, request *http.Request, context common.IRequestContext) error {
	// do nothing
	return nil
}
