package main

import (
	"log"
	"os"
)

// map levels to int for later use, like logger.Log(INFO, "%s", "logging!")
const (
	DEBUG int = iota + 1
	INFO
	WARN
	ERROR
	OFF // this level will not log anything as everything is lower
)

var (
	Std = log.New(os.Stderr, "", log.LstdFlags)

	// map to assign the level to a string, i.e: lvlMap[INFO] == "INFO"
	lvlMap = map[int]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
		OFF:   "",
	}
)

type (
	LoggerFunc func(int, string, ...interface{})
	Logger     interface {
		Log(level int, format string, a ...interface{})
	}
	WriterFunc func(p []byte) (n int, err error)
	// our structure, contains a fixed level and two hidden members, those are here
	// for internal work.
	PrefixLevelLogger struct {
		Level  int
		prefix string
		logger *log.Logger
	}
)

func (f LoggerFunc) Log(level int, format string, a ...interface{}) { f(level, format, a...) }

func NewNoopLogger() Logger { return LoggerFunc(func(_ int, _ string, _ ...interface{}) {}) }

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
