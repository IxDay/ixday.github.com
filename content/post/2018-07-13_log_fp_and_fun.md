+++
title = "Log, FP and fun in Golang"
date = 2018-07-13
url = "post/log_fp_and_fun"
+++

Golang implementation of a simple logger is subject to some criticism. First,
it does not use any extendable pattern, the only interface we can inject to
modify the behavior of the logger is a [writer interface](https://golang.org/pkg/log/#New).
Also, the lack of a proper log level control is a subject of discussion in
the community [here](https://twitter.com/bketelsen/status/820768241849077760),
[here](https://groups.google.com/forum/#!topic/golang-dev/F3l9Iz1JX4g), or
[this](https://dave.cheney.net/2015/11/05/lets-talk-about-logging) blog post.

This post is more a reflexion around APIs and tricks in Golang than a production
ready package. Take things from this post with a grain of salt as there's always
a more solutions to a given problem, especially in programming.

The final implementation is living in
[my snippet repository](https://github.com/IxDay/snippets/tree/master/golang/src/logger).

Define the basis
----------------

What do we need? A way to log at a certain level and pass some computed
values as a context. The `log.Printf` function is a good start and we will
just add a level to log to in order to be able to perform some filtering.

```go
package main

type (
	// Basic interface
	Logger interface {
		Log(level int, format string, a ...interface{})
	}

	// Convenient type to inline interface with a function
	LoggerFunc func(int, string, ...interface{})
)

func (f LoggerFunc) Log(level int, format string, a ...interface{}) { f(level, format, a...) }

func NewNoopLogger() Logger { return LoggerFunc(func(_ int, _ string, _ ...interface{}) {}) }
```

Now we'd like to be able to use the stdlib `*log.Logger` object to manage the
burden of writing to the filesystem and filter messages on a fixed mark.

```go
func BaseLogger(mark int, logger *log.Logger) Logger {

	// Use the inline function to avoid declaring unneeded structs
	return LoggerFunc(func(level int, format string, a ...interface{}) {

		// Ensure that the level is higher than the one we fixed
		if level >= mark {

			// Delegate to stdlib logger
			logger.Printf(format, a...)
		}
	})
}
```

And just a simple enumeration to make this `lvl` integer a bit more understandable.
A basic logger can also make things easier (one with flags the other without).

```go
const (
	DEBUG int = iota + 1
	INFO
	WARN
	ERROR
	OFF
)

var (
	Std = log.New(os.Stderr, "", log.LstdFlags)
	StdF = log.New(os.Stderr, "", 0)
)
```

Now we can log as we want!

```go
func main() {
	l := BaseLogger(WARN, Std)
	l.Log(DEBUG, "debug log level")
	l.Log(INFO, "info log level")
	l.Log(WARN, "warn log level")
	l.Log(ERROR, "error log level")
}
```

Going further!
--------------

Basis is here, let's see how we can go further. For instance, let's
add some prefixes:

```go
func LevelLogger(logger Logger) Logger {
	return LoggerFunc(func(level int, format string, a ...interface{}) {
		logger.Log(level, "["+LvlMap[level]+"] "+format, a...)
	})
}

func PrefixLogger(prefix string, logger Logger) Logger {
	if prefix != "" {
		prefix = prefix + ": "
	}
	return LoggerFunc(func(level int, format string, a ...interface{}) {
		logger.Log(level, prefix+format, a...)
	})
}
```

Now lets compose and print:

```go
func main() {
	l := PrefixLogger("foo", LevelLogger(BaseLogger(WARN, StdF)))

	// print nothing debug < warn
	l.Log(DEBUG, "debug log level")
	// print nothing info < warn
	l.Log(INFO, "info log level")
	// print "[WARN] foo: warn log level"
	l.Log(WARN, "warn log level")
	// print "[ERROR] foo: error log level"
	l.Log(ERROR, "error log level")
}
```

Going functional programming
----------------------------

All of those are functions, let's see what we can do to compose them and
create loggers with more features!

First, let's compose loggers

```go
// Here we take an array of callback functions, it allows us to pass wrapped
// loggers, but keep the possibility of injecting paramaters with closures
func ComposerLogger(loggers ...func(Logger) Logger) Logger {

	// Initialise with a the noop logger, you will have to provide a base logger
	// for something to happen
	base := NoopLogger()
	for _, logger := range loggers {

		// Wrap previous with new one
		base = logger(base)
	}
	return base
}
```

And now we can provide some new kinds of loggers by combining others, for instance:

```go
func PrefixLevelLogger(prefix string, mark int, logger *log.Logger) Logger {
	return ComposerLogger(
		func(_ Logger) Logger { return BaseLogger(mark, logger) },
		func(logger Logger) Logger { return LevelLogger(logger) },
		func(logger Logger) Logger { return PrefixLogger(prefix, logger) },
	)
}
```

__FUNCTIONNAL PROGRAMMING!__

We can go for a lot of different configurations, there's more in my
[snippet repo](https://github.com/IxDay/snippets/tree/master/golang/src/logger) if you need more example (there's a systemd level logger, or even a json one check it out).

Hotpatching Stdlib
------------------

Some libraries are logging directly using the `log.Printf`, `log.Println`, ...
It's possible to plug in there and here is an example:

```go

// We create a singleton of patchers
var LogPatchers = []io.Writer{}

// A bit of boilerplate
type (
	WriterFunc  func(p []byte) (n int, err error)
)
func (wf WriterFunc) Write(p []byte) (n int, err error) { return wf(p) }

// Use golang init function
func init() {
	// Hot patch log, remove flags to ensure to avoid parsing unwanted strings
	log.SetFlags(0)

	// Override the output from standard lib
	log.SetOutput(WriterFunc(func(p []byte) (int, error) {
		// Write to all our patchers
		return io.MultiWriter(LogPatchers...).Write(p)
	}))
}

func main() {
	LogPatchers = append(
		// Add our own log patcher which logs everything at the ERROR level
		// there's possibility to do some string matching if you want to
		// filter logs.
		LogPatchers, WriterFunc(func(p []byte) (int, error) {
			PrefixLevelLogger("stdlib", ERROR, Std).Log(ERROR, "%s", p)

			// Just mark the operation as a success
			return len(p), nil
		}),
	)
	log.Printf("error log level")
}
```

__Generate a valid logger from our custom implementation__

As we want to integrate with some other libs (like the `ErrorLog` from
[stdlib http package](https://golang.org/pkg/net/http/#Server), we need to provide
some utilitaries. Here is a quick example:

```go
// The only pluggable point is using a writer, so we pipe writes to a fixed level of log
func LoggerWriter(mark int, logger Logger) io.Writer {
	return WriterFunc(func(p []byte) (int, error) {
		logger.Log(mark, "%s", p)
		return len(p), nil
	})
}

// And we also provide a ready to use function, writing to a logger with no flags
func StdLogger(mark int, logger Logger) *log.Logger {
	return log.New(LoggerWriter(mark, logger), "", 0)
}
```
