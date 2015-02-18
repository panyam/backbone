package gorilla

import (
	"net/http"
	// "log"
	// "time"
)

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
}

/**
 * Middleware functions can be applied to request in a modular fashion.
 *
 * Params:
 * forward	-	Indicates whether the middleware is being applied in the forward
 * 				(request processing) or reverse (response processing) direction.
 * rw		-	The response writer for the request
 * request	-	The request being handled.
 * context	-	The request context that keeps track of state across middleware
 * 				and handlers.
 *
 * Returns:
 * A MiddlewareResult object.
 */
type MiddlewareFunc func(forward bool,
	rw http.ResponseWriter,
	request *http.Request,
	context *RequestContext) MiddlewareResult

/**
 * Stores result of a middleware function application.
 * status = 0	=>	Then middleware returned synchronously WITH A SUCCESS with
 * 					an optional result value stored (stored in the result member).
 *
 * status = -1  =>  Middleware returned synchronously WITH A FAILURE with the
 * 					optional error  in the error parameter (and any other
 * 					details in the result value)
 *
 * status = 1 	=>	Middleware is performing an asynch operation and will
 * 					eventually return.  The second return value will be a
 * 					channel of type MessageResult in which the eventual result
 * 					(success or failure) will be passed.
 */
type MiddlewareResult struct {
	/**
	 * Tells if the result is synchronous or asynchronous.
	 * If the value is synchronous then value and error are valid now.
	 * Otherwise, value is a channel of type MiddlewareResult.
	 */
	status int
	error  error
	value  interface{}
}

func MiddlewareResultSync(value interface{}, error error) MiddlewareResult {
	status := 0
	if error != nil {
		status = -1
	}
	return MiddlewareResult{status: status, value: value, error: error}
}

func MiddlewareResultAsync(channel chan MiddlewareResult) MiddlewareResult {
	return MiddlewareResult{status: 1, value: channel}
}

func (m *MiddlewareResult) IsAsync() bool {
	return m.status > 0
}

func (m *MiddlewareResult) Error() error {
	return m.error
}

func (m *MiddlewareResult) Value() interface{} {
	return m.value
}

type MiddlewareChain struct {
	middlewares     []MiddlewareFunc
	isForwardAsync  []bool
	isBackwardAsync []bool
	currIndex       int
}

func NewMiddlewareChain() *MiddlewareChain {
	out := MiddlewareChain{}
	out.middlewares = make([]MiddlewareFunc, 0, 10)
	out.isForwardAsync = make([]bool, 0, 10)
	out.isBackwardAsync = make([]bool, 0, 10)
	return &out
}

func (m *MiddlewareChain) Current() MiddlewareFunc {
	return m.middlewares[m.currIndex]
}

/**
 * Adds a new middleware to the chain.
 */
func (mc *MiddlewareChain) Add(m MiddlewareFunc, fwdAsync bool, backAsync bool) {
	mc.isForwardAsync = append(mc.isForwardAsync, fwdAsync)
	mc.isBackwardAsync = append(mc.isBackwardAsync, backAsync)
	mc.middlewares = append(mc.middlewares, m)
}

func (m *MiddlewareChain) Prev() bool {
	if m.currIndex-1 < 0 {
		return false
	}
	m.currIndex--
	return true
}

func (m *MiddlewareChain) Next() bool {
	if m.currIndex+1 >= len(m.middlewares) {
		return false
	}
	m.currIndex++
	return true
}

/**
 * Returns a final request handler function that is the result of applying all
 * the middleware in sequence and back.
 */
func (mchain *MiddlewareChain) Apply(handler RequestHandlerFunc) HttpHandlerFunc {
	// apply middleware 1
	// get its result (either sync or async)
	// if error then stop and start unwinding back the middleware chain
	return func(rw http.ResponseWriter, request *http.Request) {
		context := NewRequestContext()
		forward := true
		for {
			middleware := mchain.Current()
			result := middleware(forward, rw, request, context)

			// wait for async results
			for result.IsAsync() {
				value_chan := result.value.(chan MiddlewareResult)
				result = <-value_chan
			}

			if result.Error() != nil {
				// start going backwards
				context.AddError(result.Error())
				forward = false
			}
			if forward {
				if !mchain.Next() {
					forward = false
					// call the final handler
					if handler != nil {
						handler(rw, request, context)
					}
				}
			} else if !mchain.Prev() {
				break
			}
		}
	}
}
