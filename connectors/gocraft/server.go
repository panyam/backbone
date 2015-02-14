package gocraft

import (
	"github.com/gocraft/web"
	"net/http"
	// "fmt"
	// "strings"
)

type Context struct {
	HelloCount int
}

func (c *Context) SetHelloCount(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	c.HelloCount = 3
	next(rw, req)
}

func (c *Context) Dummy(rw web.ResponseWriter, req *web.Request) {
}

type Server struct {
}

func (s *Server) Run() {
	rootRouter := web.New(Context{}). // Create your router
						Middleware(web.LoggerMiddleware).
						Middleware(web.ShowErrorsMiddleware)

	webRouter := rootRouter.Subrouter(Context{}, "/web")
	webRouter.Middleware(StaticMiddleware("/Users/sri/projects/go/src/github.com/panyam/backbone/clients"))
	webRouter.Get("/:path:.*", (*Context).Dummy)

	http.ListenAndServe("localhost:3000", rootRouter) // Start the server!
}

func NewServer() *Server {
	return &Server{}
}
