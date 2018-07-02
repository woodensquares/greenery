// Package greenery/zapbackend is a sample zap logging back-end for greenery
// based applications.
/*
The zapbackend package implements the greenery library logging interface,
using https://github.com/uber-go/zap/ as its back-end. Compared to the
standard baselogging package, the zapbackend package supports proper
structured logging as well as ANSI colored pretty logging.

The Custom function allows logging with arbitrary zapcore.Field data, if
maximum speed and flexibility are required.
*/
package zapbackend
