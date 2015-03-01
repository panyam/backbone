package auth

import (
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

type RequestContext interface {
	AddError(err error)
	Set(key string, value interface{})
	Get(key string) interface{}
}

/**
 * Applies auth over a few methods to get a User object.
 */
func AuthMiddleware(forward bool, rw http.ResponseWriter,
	request *http.Request, context *RequestContext) MiddlewareResult {
}

func PasswordValidator() {
}

/**
 * Login based authenticator.
 * This checks that session ID and token values are valid and retrieves the user
 * based on this.
 */
type SessionValidator struct {
}

func (fm *SessionValidator) ProcessRequest(rw http.ResponseWriter, request *http.Request, context *RequestContext) MiddlewareResult {
}

func (fm *SessionValidator) ProcessResponse(rw http.ResponseWriter, request *http.Request, context *RequestContext) MiddlewareResult {
}

/**
 * Validates request based on token, secret key and signature.
 */
type SignedRequestValidator struct {
}

func (fm *SignedRequestValidator) ProcessRequest(rw http.ResponseWriter, request *http.Request, context *RequestContext) MiddlewareResult {
}

func (fm *SignedRequestValidator) ProcessResponse(rw http.ResponseWriter, request *http.Request, context *RequestContext) MiddlewareResult {
}
