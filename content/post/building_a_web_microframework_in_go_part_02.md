+++
title = "Building a web microframework in Go - Part 02"
date = 2018-05-24T10:02:07+02:00
draft = true
categories = ["Tuto"]
tags = ["golang", "dev"]
+++


__First part available [here]({{< relref "building_a_web_microframework_in_go_part_01.md" >}}).__

Previous post is focusing on Golang architecture, and how we can plug some parts
in `http.Server`. We will now go a bit deeper and starts implementing some
boilerplate abstracting a bit from stdlib primitives.

Architecture
------------

We will implement a new top level object, like `http.Server` but with our
custom logic and change a bit the API. From previous post, we defined a custom
logger. The idea here is to plug this logger in every part of `http.Server`.

```go
import (
	"bytes"
	"net/http"
	"os"
	"time"
)

const (
	DEFAULT_ADDR = "localhost:8000"
	LOG_PREFIX   = "myapp."
)

type (
	// simple composition we keep the server from stdlib and plug our logger
	Server struct {
		*http.Server
		*http.ServeMux // we also plug a mux router for now
		Logger
	}
)

// override the default method to add some insights
func (s *Server) ListenAndServe() error {

	// error handling delegation to proper logs
	s.Log(DEBUG, "Debug mode enabled")
	s.Log(INFO, "Serving on %s...", s.Server.Addr)
	if err := s.Server.ListenAndServe(); err != nil {
		s.Log(ERROR, "Failed to start server: %q", err)
		return err
	}
	return nil
}

func NewServer(options ...func(*Server)) *Server {
	router := http.NewServeMux()
	// we setup some default values, this will work out of the box when reused
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
	}

	// ensure that default can be customized
	for _, option := range options {
		option(server)
	}

	return server
}

// http.Server.ErrorLog write logs in the form: "http: some error", we don't
// need the "http:" bits, we remove it by using the callback which splits the first
// ": " we encounter
func cb(p []byte) []byte { return bytes.SplitN(p, []byte{':', ' '}, 2)[1] }

func main() {
	// create a logger
	logger := NewPrefixLevelLogger(PrefixOpt(LOG_PREFIX + "server"))

	// replace the parts of the server with our custom logger
	server := NewServer(func(s *Server) {
		s.Server.ErrorLog = logger.Logger(ERROR, cb)
		s.Logger = logger
	})

	// run the server, proper logging as been delegated to overriden method
	if err := server.ListenAndServe(); err != nil {
		os.Exit(1)
	}
}
```

Router
------

First, we'd like to use another router, we may need better performance or more
feature than what is provided by `http.ServeMux`. For this example I'd like to
use [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter),
it parses the url and provide parameters in a dedicated struct:
[see here](https://godoc.org/github.com/julienschmidt/httprouter#Handle).
To perform this, I replace the router in the structure initializer:

```go
import (
	"github.com/julienschmidt/httprouter"
)

type (
	// simple composition we keep the server from stdlib and plug our logger
	Server struct {
		*http.Server
		*httprouter.Router // replace MuxServer with httprouter Router
		Logger
	}
)

func NewServer(options ...func(*Server)) *Server {
	router := httprouter.New() // new router here
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
	}
	// ...
```

You also need this package in your `GOPATH`, retrieving the package is as simple
as `go get "github.com/julienschmidt/httprouter"`.

__Note:__ The idea here is to plug whatever fits you the most, I use `julienschmidt`
implementation, but nothing stops you from using [Gorilla's router](https://godoc.org/github.com/gorilla/mux#Router)


Handlers
--------

Now we'd like to improve the way we are dealing with http requests. Most
of the libraries out there use a context and most of the time return an error
if something goes wrong:

- Gin: [context](https://godoc.org/github.com/gin-gonic/gin#Context),
	[handler](https://godoc.org/github.com/gin-gonic/gin#HandlerFunc)
- Echo: [context](https://godoc.org/github.com/labstack/echo#Context),
	[handler](https://godoc.org/github.com/labstack/echo#HandlerFunc)
- Beego: [context](https://godoc.org/github.com/astaxie/beego/context#Context),
	[handler](https://godoc.org/github.com/astaxie/beego#FilterFunc)

Here is a way to implement this:

```go
import (
	"bytes"
	"io"
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
)

//  Wrap the httprouter Handle method
func (s *Server) Handle(method string, path string, handler Handler) {
	wrapped := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		// build the context
		context := &Context{w, r, p, s.Logger}

		// pass it to our handler
		if err := handler(context); err != nil {
			// here we will handle the the error
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
```

Errors
------

With the new way of handling requests we are now returning errors, the last
step is to plug here and provide more flexibility.

```go
import (
	"fmt"
)

// define the default way to handle error
func DefaultErrorHandler(err error, ctx *Context) {
	ctx.Log(ERROR, "Failed serving %q: %q", ctx.Request.URL, err)
	ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
	ctx.ResponseWriter.Write([]byte(fmt.Sprintf("Internal server error: %q", err)))
}

// now lets plug it to our handle function
func (s *Server) Handle(method string, path string, handle Handle) {
	// ...
		// pass it to our handler
		if err := handle(context); err != nil {
			// if something bad happen
			s.ErrorHandler(err, context)
		}
	// ...
}

// and make the handler configurable at server struct level
type (
	ErrorHandler func(error, *Context)
	Server struct {
		*http.Server
		*httprouter.Router
		Logger
		ErrorHandler ErrorHandler
	}
)

// provide a default working implementation
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
	// ...
}
```

__TODO:__ add something here

A basic code example is available [in the blog repo](https://github.com/IxDay/ixday.github.com/tree/source/content/code/microframework_in_go/part_02).

