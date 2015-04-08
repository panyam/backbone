package middleware

import (
	"errors"
	"fmt"
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
	log.Println("Validating request: ", request)
	if am.Validators != nil {
		for _, validator := range am.Validators {
			user := validator.ValidateRequest(request, context)
			if user != nil {
				context.Set("user", user)
				return nil
			}
		}
	}

	// none of the validators matched so throw a 401 unauthorized response
	rw.WriteHeader(http.StatusUnauthorized)
	log.Println("No validator matched request, request: ", request)
	msg := fmt.Sprintf("%d Unauthorized", http.StatusUnauthorized)
	err := errors.New(msg)
	// utils.SendJsonError(rw, map[string]interface{}{"error": msg}, http.StatusUnauthorized)
	return err
}

func (validator *AuthMiddleware) ProcessResponse(rw http.ResponseWriter, request *http.Request, context common.IRequestContext) error {
	// do nothing
	return nil
}
