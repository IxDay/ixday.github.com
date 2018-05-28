+++
title = "Building a web microframework in Go - Part 01"
date = 2018-05-11T20:37:48+02:00
draft = true
categories = ["Tuto"]
tags = ["golang", "dev"]
+++

__Should I use a web framework when building a website or tie to the stdlib
primitives?__

This is an heavy debated question in the Golang community, and both
have pro and cons.

I built a various set of services using both approach and now I realize that
building your own wrappers around stdlib is not that hard. So, in this serie of
posts I will show you how to do it:

- __Part One:__ write clean and extensible golang structures, we will start with
	a level based prefix logger.
- __Part Two:__ customize the handler, intercept stdlib `http.Handler` calls, create
	a context with our own needs and inject it inside our custom handler interface.
	This article will also contains how to tackle error handling, and custom
	router.
- __Part Three:__ write middleware and using a rendering engine.

__Note:__ I will write a lot of code in the posts and it will be hard to
follow up due to the fragmented parts. For this reason I put the code in
this blog repository to test it. You can find those resources
[here](https://github.com/IxDay/ixday.github.com/tree/source/content/code/microframework_in_go).

State of the ecosystem
----------------------

Currently Golang ecosystem is splitted between those two approach. In order to
build something it's important to get a tour of what is already existing.

First, you can do without, 5mn on Google about golang frameworks will _propulse_ you in
[threads were people will advocate for only using stdlib](https://news.ycombinator.com/item?id=11427542).
This is true if you are doing a simple JSON REST API but if you are going for
a full service things will become harder pretty fast.
The `Handler` interface from stdlib is a good abstraction
but it leaves you with only a few number of helpers an almost no guidance.

Then, there's frameworks, and you have a lot of solutions out
[there](https://blog.usejournal.com/top-6-web-frameworks-for-go-as-of-2017-23270e059c4b).
Currently, no clear leader as emerged yet, and the only standard way of plugging here is to use
[the stdlib handler interface](https://golang.org/pkg/net/http/#Handler).
However, those solutions usually provide convenient utilities to start writing
web services, and we will try to mimic some of them.

Both approaches have pros and cons:

- Framework: default architecture, faster start. Problem appears when you
want to go out of the common use cases.
- Stdlib: simple primitives which give full control. Code can get messy if
you do not think of the architecture ahead.

Architecture and write some Golang
----------------------------------

Before starting, I'd like to share some resources which helped me a lot in writing
Go. I am extensively writing Golang for more than two years now, I've read:
[The go programming language book](http://www.gopl.io/) but nothing improved my
Golang programming skills as those two resources:

- [Youtube video around the interface paradigm in golang](https://www.youtube.com/watch?v=xyDkyFjzFVc)
- [Dave Cheney post around functional options](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)

First module
------------

__This example is here to think around interfaces and Golang APIs.
It is not meant to be a production ready solution.
You can (and maybe should) use more advanced modules to handle logs in your project.__

Golang implementation of a simple logger is subject to some criticism. First,
it does not use any extendable pattern, the only interface we can inject to
modify the behavior of the logger is a [writer interface](https://golang.org/pkg/log/#New).
Also, the lack of a proper log level control is a subject of discussion in
the community [here](https://twitter.com/bketelsen/status/820768241849077760),
[here](https://groups.google.com/forum/#!topic/golang-dev/F3l9Iz1JX4g), or
[this](https://dave.cheney.net/2015/11/05/lets-talk-about-logging) blog post.

Let's define what we need: a way to log at a certain level
(we will also use what was described in the youtube video from
[previous section](#architecture-and-write-some-golang)).

```go
package main

type (
	LoggerFunc func(int, string, ...interface{})
	Logger     interface {
		Log(level int, format string, a ...interface{})
	}
)

func (f LoggerFunc) Log(level int, format string, a ...interface{}) { f(level, format, a...) }

func NewNoopLogger() Logger { return LoggerFunc(func(_ int, _ string, _ ...interface{}) {}) }
```

It defines a simple Logger interface with a `Log` function, this function
is using the same signature as `fmt.Printf` and `log.Printf` with an additional
argument: a level at which we want to log.
We also define a `NoopLogger` in order to discard logs.

Now, let's implement our prefix level logger:

```go
import (
	"log"
)

// map levels to int for later use, like logger.Log(INFO, "%s", "logging!")
const (
	DEBUG int = iota + 1
	INFO
	WARN
	ERROR
	OFF  // this level will not log anything as everything is lower
)

// map to assign the level to a string, i.e: lvlMap[INFO] == "INFO"
var (
	lvlMap = map[int]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
		OFF:   "",
	}
)

// our structure, contains a fixed level and two hidden members, those are here
// for internal work.
type (
	PrefixLevelLogger struct {
		Level  int
		prefix string
		logger *log.Logger
	}
)

// as prefix need to be reworked before being set we make it inaccessible
func (pll *PrefixLevelLogger) Prefix() string { return pll.prefix[:len(pll.prefix)-2] }
func (pll *PrefixLevelLogger) SetPrefix(p string) {
	if p != "" {
		pll.prefix = p + ": "
	}
}

// logger will need to be a bit more reworked by the getter, we will show this later
func (pll *PrefixLevelLogger) SetLogger(l *log.Logger) { pll.logger = l }

// most important function, it is responsible for logging lines, if level is sufficient
func (pll *PrefixLevelLogger) Log(level int, format string, a ...interface{}) {
	if level >= pll.Level {
		a = append([]interface{}{lvlMap[level], pll.prefix}, a...)
		pll.logger.Printf("[%s] %s"+format, a...)
	}
}
```

First bits are here, this should be able to log, we will now use what Dave Cheney
taught us to build a functionnal API.

```go
import (
	"os"
)

var (
	Std = log.New(os.Stderr, "", log.LstdFlags)
)

// Define a new prefix logger, with a set of options to configure it
func NewPrefixLevelLogger(options ...func(*PrefixLevelLogger)) *PrefixLevelLogger {

	// build a default working logger, logging at INFO level to Stderr with date and time
	pll := &PrefixLevelLogger{INFO, "", Std}

	// use options to customize our logger
	for _, option := range options {
		option(pll)
	}
	return pll
}

// Define some convenient wrappers for setting up options
func PrefixOpt(prefix string) func(*PrefixLevelLogger) {
	return func(pll *PrefixLevelLogger) { pll.SetPrefix(prefix) }
}

func LevelOpt(level int) func(*PrefixLevelLogger) {
	return func(pll *PrefixLevelLogger) { pll.Level = level }
}

func LoggerOpt(logger *log.Logger) func(*PrefixLevelLogger) {
	return func(pll *PrefixLevelLogger) { pll.SetLogger(logger) }
}

// This one allows us to clone a previously created logger to avoid resetting one from scratch
func CloneOpt(pll *PrefixLevelLogger) func(*PrefixLevelLogger) {
	return func(_pll *PrefixLevelLogger) {
		_pll.Level, _pll.prefix, _pll.logger = pll.Level, pll.prefix, pll.logger
	}
}
```

Okay now just a bit more boilerplate, to make it even more extensible. We will
provide a `Clone` and a `Logger` method. The first one will provide a cleaner API
the second will allow to inject our logger in external services relying on
`log.Logger`. Those are specifics to our custom structure.

We could have added it to the `Logger` interface but it may have cluttered it.
For the moment we will stay with it, but later we can create a `Cloner` interface and add
a `CloneLogger` composition. A good example of interface composition can
be found in [the io package](https://golang.org/pkg/io/#ReadCloser).

```go
type (
	WriterFunc func(p []byte) (n int, err error)
)

func (wf WriterFunc) Write(p []byte) (n int, err error) { return wf(p) }

// signature is the same as a new logger, but populated with the current one fields
func (pll *PrefixLevelLogger) Clone(options ...func(*PrefixLevelLogger)) *PrefixLevelLogger {
	return NewPrefixLevelLogger(
		append([]func(*PrefixLevelLogger){CloneOpt(pll)}, options...)...,
	)
}

// we want to provide a legacy *log.Logger.
// the idea here is to create one which will log at a specific level
// third arg allows you to parse what come through in case of edit, see how it is used with http.Server
func (pll *PrefixLevelLogger) Logger(level int, cb func(p []byte) []byte) *log.Logger {
	// create a *log.Log with our writer wrapper. we then pass it to our underlying
	// logger at the specified level
	return log.New(WriterFunc(func(p []byte) (int, error) {
		pll.Log(level, "%s", cb(p))
		return len(p), nil
	}), "", 0)
}
```

This example also shows how to use Golang attributes visibility. The non public
members require some rework before being set or handed to the user. That's
why we are writing getters and setters and put the underlying reference as
private.

Use the logger Luke!
--------------------

Lets write a logger for `http.Server`.

```go
package main

import (
	"bytes"
	"net/http"
	"os"
)

// lets define a faulty handler to see the logger in action
func handler(w http.ResponseWriter, r *http.Request) {
	// a second call to WriteHeader trigger an error on http.Server.ErrorLog
	w.WriteHeader(http.StatusOK)
	w.WriteHeader(http.StatusInternalServerError)
}

// http.Server.ErrorLog write logs in the form: "http: some error", we don't
// need the "http:" bits, we remove it by using the callback which splits the first
// ": " we encounter
func cb(p []byte) []byte { return bytes.SplitN(p, []byte{':', ' '}, 2)[1] }

func main() {
	// create a simple logger to stderr, with "main" as a prefix
	logger := NewPrefixLevelLogger(PrefixOpt("main"))

	// now lets create a server, which will log errors at the ERROR level with another prefix
	server := &http.Server{
		Addr:     "localhost:8000",
		Handler:  http.HandlerFunc(handler),
		ErrorLog: logger.Clone(PrefixOpt("server")).Logger(ERROR, cb), // conveniently chain functions
	}
	logger.Log(INFO, "Start serving requests...")
	if err := server.ListenAndServe(); err != nil {
		logger.Log(ERROR, "Something bad happened trying to serve requests: %q", err)
		// log.Fatal does not exist anymore, we need to exit with error code
		// we can also do:
		// logger.Logger(ERROR, func([]byte){}).Fatalf("Something bad happened trying to serve requests: %q", err)
		os.Exit(1)
	}
}
```

Resources
---------

Books:

- [The go programming language book](http://www.gopl.io/)

Youtube links:

- [Youtube video around the interface paradigm in golang](https://www.youtube.com/watch?v=xyDkyFjzFVc)
- [Dave Cheney post around functional options](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)

Code:

- [Golang stdlib http](https://golang.org/pkg/net/http/)
- [This blog post code](https://github.com/IxDay/ixday.github.com/tree/source/content/code/microframework_in_go/part_01). There should be no further config required than a
go compiler and running `go run *.go`.

__Have fun!__
