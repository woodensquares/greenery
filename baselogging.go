package greenery

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Super basic logging to have a default for the library, will simply output
// everything to stderr and trace to stdout as-is, no support for structured
// logging.
const (
	debugLevel = iota
	infoLevel
	warnLevel
	errorLevel
)

// These are not public because we want users to be able to do whatever they
// want in terms of what they want to support and how to display it, so each
// logger defines their own.
const (
	genericType = iota

	durationType
	intType
	stringType
	timeType
)

type baseLogger struct {
	w      io.Writer
	level  uint8
	prefix string
}

func (l *baseLogger) DebugSkip(i int, s string) {
	// Basic logging does not care about skipping since we are not printing
	// where the log call is from.
	l.DebugStructured(s)
}

func getValues(lf []LogField) []string {
	zf := make([]string, len(lf))
	for i, v := range lf {
		switch v.Type {
		case durationType:
			zf[i] = fmt.Sprintf("%s: %v", v.Key, v.Duration)
		case timeType:
			zf[i] = fmt.Sprintf("%s: %v", v.Key, v.Time)
		case intType:
			zf[i] = fmt.Sprintf("%s: %d", v.Key, v.Integer)
		case stringType:
			zf[i] = fmt.Sprintf("%s: %s", v.Key, v.String)
		case genericType:
			zf[i] = fmt.Sprintf("%s: %v", v.Key, v.Generic)
		default:
			// Should not happen
			zf[i] = fmt.Sprintf("%s: Unsupported log type %v", v.Key, v.Type)
		}
	}

	return zf
}

// doPrint is a basic printf-er converting the various "structured" fields to
// space delimited strings.
func (l *baseLogger) doPrint(s string, lf ...LogField) {
	if len(lf) == 0 {
		fmt.Fprintf(l.w, "%s%s\n", l.prefix, s)
	} else {
		fmt.Fprintf(l.w, "%s%s %s\n", l.prefix, s, strings.Join(getValues(lf), " "))
	}
}

// No custom logging for the default base logger
func (l *baseLogger) Custom(s string, lv int, data interface{}) {}

func (l *baseLogger) DebugStructured(s string, lf ...LogField) {
	if l.level != debugLevel {
		return
	}
	l.doPrint(s, lf...)
}

func (l *baseLogger) InfoStructured(s string, lf ...LogField) {
	if l.level > debugLevel {
		return
	}
	l.doPrint(s, lf...)
}

func (l *baseLogger) WarnStructured(s string, lf ...LogField) {
	if l.level > infoLevel {
		return
	}
	l.doPrint(s, lf...)
}

func (l *baseLogger) ErrorStructured(s string, lf ...LogField) {
	l.doPrint(s, lf...)
}

// LogString returns a field representing a string value
func (l *baseLogger) LogString(name string, value string) LogField {
	return LogField{Key: name, Type: stringType, String: value}
}

// LogInteger returns a field representing an integer value
func (l *baseLogger) LogInteger(name string, value int) LogField {
	return LogField{Key: name, Type: intType, Integer: value}
}

// LogTime returns a field representing a time.Time value
func (l *baseLogger) LogTime(name string, value time.Time) LogField {
	return LogField{Key: name, Type: timeType, Time: value}
}

// LogDuration returns a field representing a time.Duration value
func (l *baseLogger) LogDuration(name string, value time.Duration) LogField {
	return LogField{Key: name, Type: durationType, Duration: value}
}

// LogGeneric returns a field representing a generic value
func (l *baseLogger) LogGeneric(name string, value interface{}) LogField {
	return LogField{Key: name, Type: genericType, Generic: value}
}

// Sync is called on application exit typically and would sync the log to
// storage if needed, no need to do anything here as stdout/err are not
// buffered by default.
func (l *baseLogger) Sync() error {
	return nil
}

// --------------------------------------------------------------------------

func logValueFromString(e string) uint8 {
	switch strings.ToLower(e) {
	case "error":
		return errorLevel
	case "warn":
		return warnLevel
	case "info":
		return infoLevel
	case "debug":
		return debugLevel
	}

	// This really shouldn't happen, the validation is done at enum
	// setting time
	panic(fmt.Sprintf("Invalid log level %s, only error, warn, info and debug are supported", e))
}

// BaseTraceLogger represents a base logger used for Tracing
func BaseTraceLogger(cfg Config) Logger {
	return &baseLogger{w: os.Stdout, level: debugLevel, prefix: "TRACE: "}
}

// BasePrettyLogger represents a base logger used for pretty logging
func BasePrettyLogger(cfg Config, l string, w io.Writer) Logger {
	return &baseLogger{w: w, level: logValueFromString(l)}
}

// BaseStructuredLogger represents a base logger used for structured logging,
// although for base logs this is not supported and it will behave the same
// way as the pretty logger for the time being
func BaseStructuredLogger(cfg Config, l string, w io.Writer) Logger {
	return BasePrettyLogger(cfg, l, w)
}
