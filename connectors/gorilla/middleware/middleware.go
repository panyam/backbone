package middleware

import (
	. "github.com/panyam/relay/connectors/gorilla/common"
	"log"
	"net/http"
)

/**
 * Middleware functions can be applied to request in a modular fashion.
 *
 * Params:
 * rw		-	The response writer for the request
 * request	-	The request being handled.
 * context	-	The request context that keeps track of state across middleware
 * 				and handlers.
 *
 * Returns:
 * A error if any
 */
type MiddlewareFunc func(rw http.ResponseWriter, request *http.Request, context IRequestContext) error

type Middleware interface {
	ProcessRequest(rw http.ResponseWriter, request *http.Request, context IRequestContext) error
	ProcessResponse(rw http.ResponseWriter, request *http.Request, context IRequestContext) error
}

type FunctionMiddleware struct {
	requestMiddleware  MiddlewareFunc
	responseMiddleware MiddlewareFunc
}

func (fm *FunctionMiddleware) ProcessRequest(rw http.ResponseWriter, request *http.Request, context IRequestContext) error {
	if fm.requestMiddleware != nil {
		return fm.requestMiddleware(rw, request, context)
	}
	return nil
}

func (fm *FunctionMiddleware) ProcessResponse(rw http.ResponseWriter, request *http.Request, context IRequestContext) error {
	if fm.responseMiddleware != nil {
		return fm.responseMiddleware(rw, request, context)
	}
	return nil
}

type MiddlewareChain struct {
	middlewares []Middleware
	currIndex   int
}

func NewMiddlewareChain() *MiddlewareChain {
	out := MiddlewareChain{}
	out.middlewares = make([]Middleware, 0, 10)
	return &out
}

func (m *MiddlewareChain) Current() Middleware {
	if m.currIndex >= len(m.middlewares) {
		return nil
	}
	return m.middlewares[m.currIndex]
}

/**
 * Adds a new middleware to the chain.
 */
func (mc *MiddlewareChain) AddMiddleware(m Middleware) {
	mc.middlewares = append(mc.middlewares, m)
}

func (mc *MiddlewareChain) AddRequestMiddleware(requestMiddleware MiddlewareFunc) {
	mc.AddMiddleware(&FunctionMiddleware{requestMiddleware, nil})
}

func (mc *MiddlewareChain) AddResponseMiddleware(responseMiddleware MiddlewareFunc) {
	mc.AddMiddleware(&FunctionMiddleware{responseMiddleware, nil})
}

func (mc *MiddlewareChain) Add(requestMiddleware MiddlewareFunc, responseMiddleware MiddlewareFunc) {
	mc.AddMiddleware(&FunctionMiddleware{requestMiddleware, responseMiddleware})
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
	if mchain.Current() == nil {
		return func(rw http.ResponseWriter, request *http.Request) {
			handler(rw, request, nil)
		}
	}
	return func(rw http.ResponseWriter, request *http.Request) {
		context := NewRequestContext()
		forward := true
		for {
			middleware := mchain.Current()
			var err error = nil
			if forward {
				err = middleware.ProcessRequest(rw, request, context)
			} else {
				err = middleware.ProcessResponse(rw, request, context)
			}

			if err != nil {
				// start going backwards
				log.Println("Adding error..., forward: ", forward)
				context.AddError(err)
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
