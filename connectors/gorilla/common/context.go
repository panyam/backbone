package common

import "log"

type IRequestContext interface {
	AddError(err error)
	Set(key string, value interface{})
	Get(key string) interface{}
	Errors() []error
}

/**
 * Request context stores all data related to a request handling sessions.
 */
type RequestContext struct {
	errors []error
	data   map[string]interface{}
}

func NewRequestContext() *RequestContext {
	rc := RequestContext{data: make(map[string]interface{})}
	return &rc
}

func (rc *RequestContext) AddError(err error) {
	rc.errors = append(rc.errors, err)
	log.Println("Added error: ", err)
}

func (rc *RequestContext) Set(key string, value interface{}) {
	rc.data[key] = value
}

func (rc *RequestContext) Get(key string) interface{} {
	return rc.data[key]
}

func (rc *RequestContext) Errors() []error {
	return rc.errors
}
