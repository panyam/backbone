package common

import (
	"net/http"
)

type HttpHandlerFunc func(http.ResponseWriter, *http.Request)
type RequestHandlerFunc func(http.ResponseWriter, *http.Request, *RequestContext)
type ServiceRequestMaker func(*http.Request) (interface{}, error)
type ServiceResponseMaker func(http.ResponseWriter, *http.Request, interface{}, error)
