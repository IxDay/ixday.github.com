package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

const (
	DEFAULT_ADDR = "localhost:8000"
)

type (
	// here is the function which will be used to handle requests
	Handler func(*Context) error

	// custom context, I put here what I need. You can also push your database client for example
	Context struct {
		http.ResponseWriter
		*http.Request
		httprouter.Params
		Logger
	}

	ErrorHandler func(error, *Context)
	// simple composition we keep the server from stdlib and plug our logger
	Server struct {
		*http.Server
		*httprouter.Router // replace MuxServer with httprouter Router
		Logger
		ErrorHandler ErrorHandler
	}
)

//  Wrap the httprouter Handle method
func (s *Server) Handle(method string, path string, handler Handler) {
	wrapped := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		// build the context
		context := &Context{w, r, p, s.Logger}

		// pass it to our handler
		if err := handler(context); err != nil {
			// if something bad happen
			s.ErrorHandler(err, context)
		}
	}
	// plug it to the server router
	s.Router.Handle(method, path, wrapped)
}

// lets add some handy helpers
func (s *Server) GET(path string, handler Handler)  { s.Handle("GET", path, handler) }
func (s *Server) POST(path string, handler Handler) { s.Handle("POST", path, handler) }
func (s *Server) PUT(path string, handler Handler)  { s.Handle("PUT", path, handler) }
func (s *Server) HEAD(path string, handler Handler) { s.Handle("HEAD", path, handler) }

// override the default method to add some insights
func (s *Server) ListenAndServe() error {
	s.Log(DEBUG, "Debug mode enabled")
	s.Log(INFO, "Serving on %s...", s.Server.Addr)
	if err := s.Server.ListenAndServe(); err != nil {
		s.Log(ERROR, "Failed to start server: %q", err)
		return err
	}
	return nil
}

// define the default way to handle error
func DefaultErrorHandler(err error, ctx *Context) {
	ctx.Log(ERROR, "Failed serving %q: %q", ctx.Request.URL, err)
	ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
	ctx.ResponseWriter.Write([]byte(fmt.Sprintf("Internal server error: %q", err)))
}

func NewServer(options ...func(*Server)) *Server {
	router := httprouter.New()
	server := &Server{
		&http.Server{
			Addr:         DEFAULT_ADDR,
			Handler:      router,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		router,
		NewNoopLogger(),
		DefaultErrorHandler,
	}
	for _, option := range options {
		option(server)
	}

	return server
}
