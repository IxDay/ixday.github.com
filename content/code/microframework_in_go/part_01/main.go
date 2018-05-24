package main

import (
	"bytes"
	"net/http"
	"os"
)

type (
	HandleFunc func(http.ResponseWriter, *http.Request) // let's take advantage of what we learnt
)

// implement handler interface
func (hf HandleFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) { hf(w, r) }

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
		Handler:  HandleFunc(handler),
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
