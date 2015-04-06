package middleware

import (
	"github.com/panyam/relay/connectors/gorilla/common"
	// "github.com/gorilla/sessions"
	"net/http"
)

func PasswordValidator() {
}

/**
 * Login based authenticator.
 * This checks that session ID and token values are valid and retrieves the user
 * based on this.
 */
type SessionValidator struct {
}

func (fm *SessionValidator) ProcessRequest(rw http.ResponseWriter, request *http.Request, context common.IRequestContext) error {
	return nil
}

func (fm *SessionValidator) ProcessResponse(rw http.ResponseWriter, request *http.Request, context common.IRequestContext) error {
	return nil
}

/**
 * Validates request based on token, secret key and signature.
 */
type SignedRequestValidator struct {
}

func (fm *SignedRequestValidator) ProcessRequest(rw http.ResponseWriter, request *http.Request, context common.IRequestContext) error {
	return nil
}

func (fm *SignedRequestValidator) ProcessResponse(rw http.ResponseWriter, request *http.Request, context common.IRequestContext) error {
	return nil
}
