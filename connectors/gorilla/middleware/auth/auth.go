package middleware

import (
	"github.com/panyam/backbone/models"
	"log"
	"net/http"
)

type RequestContext interface {
	AddError(err error)
	Set(key string, value interface{})
	Get(key string) interface{}
}

type AuthValidator interface {
	ValidateRequest(request *http.Request, context *RequestContext) bool
}

type AuthMiddleware struct {
	Validators []AuthValidator
}

/**
 * Applies auth over a few methods to get a User object.
 */
func AuthMiddleware(forward bool, rw http.ResponseWriter,
	request *http.Request, context *RequestContext) *models.User {
}

func (am *AuthMiddleware) ProcessRequest(rw http.ResponseWriter, request *http.Request, context *RequestContext) MiddlewareResult {
	for index, validator := range am.validators {
		user := validator.ValidateRequest(request, context)
		if user != nil {
			context.Set("user", user)
			return NewMiddlewareResult(nil, nil)
		}
	}

	// none of the validators matched so throw a 401 unauthorized response
	log.Println("Auth Failed")
	rw.WriteHeader(http.StatusUnauthorized)
	return NewMiddlewareResult(nil, errors.New("Unauthorized"))
}

func (validator *AuthMiddleware) ProcessResponse(rw http.ResponseWriter, request *http.Request, context *RequestContext) MiddlewareResult {
	// do nothing
}
