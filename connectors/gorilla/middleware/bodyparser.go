package middleware

import (
	"github.com/panyam/backbone/connectors/gorilla/common"
	"net/http"
)

/**
 * Parses the request initially so subsequent middleware dont have to repeat
 * this.
 */
func BodyParserMiddleware(rw http.ResponseWriter, request *http.Request, context *common.RequestContext) MiddlewareResult {
	request.ParseForm()
	return NewMiddlewareResult(nil, nil)
}
