package main

import (
	"bytes"
	"fmt"
	"os"
)

const (
	LOG_PREFIX = "myapp"
)

// test what we just implemented by raising a status internal server error
func raiseError(c *Context) error {
	return fmt.Errorf("Should raise a 500!")
}

// http.Server.ErrorLog write logs in the form: "http: some error", we don't
// need the "http:" bits, we remove it by using the callback which splits the first
// ": " we encounter
func cb(p []byte) []byte { return bytes.SplitN(p, []byte{':', ' '}, 2)[1] }

func main() {
	// create a logger
	logger := NewPrefixLevelLogger(PrefixOpt(LOG_PREFIX + "server"))

	server := NewServer(func(s *Server) {
		s.Server.ErrorLog = logger.Logger(ERROR, cb)
		s.Logger = logger
	})
	// lets plug the handler
	server.GET("/", raiseError)
	if err := server.ListenAndServe(); err != nil {
		os.Exit(1)
	}
}
