package middleware

import (
	"net/http"
)

/**
 * Parses the request initially so subsequent middleware dont have to repeat
 * this.
 */
func BodyParserMiddleware(forward bool, rw http.ResponseWriter,
	request *http.Request, context *RequestContext) MiddlewareResult {
	request.ParseForm()
	return NewMiddlewareResult(nil, nil)
}
