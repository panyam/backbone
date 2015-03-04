package common

import (
	"net/http"
)

type HttpHandlerFunc func(http.ResponseWriter, *http.Request)
type RequestHandlerFunc func(http.ResponseWriter, *http.Request, *RequestContext)
