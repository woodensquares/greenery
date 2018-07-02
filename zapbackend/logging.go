package zapbackend

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/woodensquares/greenery"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	genericType = iota

	durationType
	intType
	stringType
	timeType
)

// The available logging levels for this backend are as follows and will map
// to the Debug, Info, Warn and Error zap methods.
const (
	Debug = iota
	Info
	Warn
	Error
)

type zapLogger struct {
	zl *zap.Logger
}

func (l *zapLogger) DebugSkip(i int, s string) {
	l.zl.WithOptions(zap.AddCallerSkip(i)).Debug(s)
}

// getValues obv allocates an intermediate object used for compatibility,
// different struct from what is created to pass to zap, if users want to
// avoid this they can of course use the logger directly, just create their
// own makelogger functions and use the zap logging calls directly in their
// code.
func getValues(lf []greenery.LogField) []zapcore.Field {
	zf := make([]zapcore.Field, len(lf))
	for i, v := range lf {
		switch v.Type {
		case durationType:
			zf[i] = zap.Duration(v.Key, v.Duration)
		case timeType:
			zf[i] = zap.Time(v.Key, v.Time)
		case intType:
			zf[i] = zap.Int(v.Key, v.Integer)
		case stringType:
			zf[i] = zap.String(v.Key, v.String)
		case genericType:
			zf[i] = zap.Any(v.Key, v.Generic)
		default:
			// Should not happen
			zf[i] = zap.String(v.Key, fmt.Sprintf("Unsupported log type %v", v.Type))
		}
	}

	return zf
}

func (l *zapLogger) Custom(s string, lv int, lc interface{}) {
	lz, ok := lc.([]zapcore.Field)
	if !ok {
		l.zl.Error("Invalid logging data", zap.String("data", spew.Sdump(lc)))
		return
	}

	switch lv {
	case Debug:
		l.zl.Debug(s, lz...)
	case Info:
		l.zl.Info(s, lz...)
	case Warn:
		l.zl.Warn(s, lz...)
	case Error:
		l.zl.Error(s, lz...)
	default:
		l.zl.Error(fmt.Sprintf("Invalid logging level %d for message %s", lv, s), lz...)
	}
}

func (l *zapLogger) DebugStructured(s string, lf ...greenery.LogField) {
	if len(lf) == 0 {
		l.zl.Debug(s)
	} else {
		l.zl.Debug(s, getValues(lf)...)
	}
}

func (l *zapLogger) InfoStructured(s string, lf ...greenery.LogField) {
	if len(lf) == 0 {
		l.zl.Info(s)
	} else {
		l.zl.Info(s, getValues(lf)...)
	}
}

func (l *zapLogger) WarnStructured(s string, lf ...greenery.LogField) {
	if len(lf) == 0 {
		l.zl.Warn(s)
	} else {
		l.zl.Warn(s, getValues(lf)...)
	}
}

func (l *zapLogger) ErrorStructured(s string, lf ...greenery.LogField) {
	if len(lf) == 0 {
		l.zl.Error(s)
	} else {
		l.zl.Error(s, getValues(lf)...)
	}
}

// LogString returns a greenery LogField for a string logging value
func (l *zapLogger) LogString(name string, value string) greenery.LogField {
	return greenery.LogField{Key: name, Type: stringType, String: value}
}

// LogInteger returns a greenery LogField for an int logging value
func (l *zapLogger) LogInteger(name string, value int) greenery.LogField {
	return greenery.LogField{Key: name, Type: intType, Integer: value}
}

// LogTime returns a greenery LogField for a time.Time logging value
func (l *zapLogger) LogTime(name string, value time.Time) greenery.LogField {
	return greenery.LogField{Key: name, Type: timeType, Time: value}
}

// LogDuration is returns a greenery LogField for a time.Duration logging value
func (l *zapLogger) LogDuration(name string, value time.Duration) greenery.LogField {
	return greenery.LogField{Key: name, Type: durationType, Duration: value}
}

// LogGeneric is returns a greenery LogField for a generic logging value
func (l *zapLogger) LogGeneric(name string, value interface{}) greenery.LogField {
	return greenery.LogField{Key: name, Type: genericType, Generic: value}
}

// Sync will be used to sync the zap log
func (l *zapLogger) Sync() error {
	return l.zl.Sync()
}

// --------------------------------------------------------------------------

// Set will validate for supported log levels and set it
func logValueFromString(e string) zapcore.Level {
	switch strings.ToLower(e) {
	case "error":
		return zapcore.ErrorLevel
	case "warn":
		return zapcore.WarnLevel
	case "info":
		return zapcore.InfoLevel
	case "debug":
		return zapcore.DebugLevel
	}

	// This really shouldn't happen, the validation is done at enum
	// setting time
	panic(fmt.Sprintf("Invalid log level %s, only error, warn, info and debug are supported", e))
}

// TraceLogger returns a greenery Logger to be used for tracing
func TraceLogger(cfg greenery.Config) greenery.Logger {
	return &zapLogger{zl: zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
				TimeKey:    "T",
				LevelKey:   "L",
				NameKey:    "N",
				CallerKey:  "C",
				MessageKey: "M",
				LineEnding: zapcore.DefaultLineEnding,
				EncodeLevel: func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
					// Zap does not support trace as a level, let's fake it in
					// our console output as a green TRACE
					enc.AppendString(fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(32), "TRACE"))
				},
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			}),
			zap.CombineWriteSyncers(zapcore.AddSync(os.Stdout)),
			zapcore.DebugLevel),
		zap.AddCallerSkip(2),
		zap.Development(),
		zap.AddCaller())}
}

// PrettyLogger returns a greenery Logger to be used for pretty logging
func PrettyLogger(cfg greenery.Config, l string, w io.Writer) greenery.Logger {
	return &zapLogger{zl: zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
				TimeKey:        "T",
				LevelKey:       "L",
				NameKey:        "N",
				CallerKey:      "C",
				MessageKey:     "M",
				StacktraceKey:  "S",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalColorLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			}),
			zap.CombineWriteSyncers(zapcore.AddSync(w)),
			zapcore.Level(logValueFromString(l))),
		zap.AddCallerSkip(2))}
}

// StructuredLogger  returns a greenery Logger to be used for structured logging
func StructuredLogger(cfg greenery.Config, l string, w io.Writer) greenery.Logger {
	return &zapLogger{zl: zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zapcore.EncoderConfig{
				TimeKey:        "ts",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "msg",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.EpochTimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			}),
			zap.CombineWriteSyncers(zapcore.AddSync(w)),
			zapcore.Level(logValueFromString(l))),
		zap.AddCallerSkip(2))}
}
