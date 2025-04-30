package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
)

// Interface -.
type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

// Logger -.
type logger struct {
	logger *zerolog.Logger
}

var _ Interface = (*logger)(nil)

// New -.
func New(level string) *logger {
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)
	skipFrameCount := 3
	z := zerolog.New(os.Stdout).
		Output(zerolog.ConsoleWriter{Out: os.Stderr, PartsOrder: []string{"level", "caller", "message"}}).
		With().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
		Logger()
	return &logger{
		logger: &z,
	}
}

func (l *logger) formatMessage(message any) string {
	switch t := message.(type) {
	case error:
		return t.Error()
	case string:
		return t
	default:
		return fmt.Sprintf("Unknown type %v", message)
	}
}

// Debug -.
func (l *logger) Debug(message interface{}, args ...interface{}) {
	mf := l.formatMessage(message)
	l.log(l.logger.Debug(), mf, args...)
}

// Info -.
func (l *logger) Info(message string, args ...interface{}) {
	mf := l.formatMessage(message)
	l.log(l.logger.Info(), mf, args...)
}

// Warn -.
func (l *logger) Warn(message string, args ...interface{}) {
	mf := l.formatMessage(message)
	l.log(l.logger.Warn(), mf, args...)
}

// Error -.
func (l *logger) Error(message interface{}, args ...interface{}) {
	mf := l.formatMessage(message)
	l.log(l.logger.Error(), mf, args...)
}

// Fatal -.
func (l *logger) Fatal(message interface{}, args ...interface{}) {
	mf := l.formatMessage(message)
	l.log(l.logger.Fatal(), mf, args...)

	os.Exit(1)
}

func (l *logger) log(e *zerolog.Event, m string, args ...interface{}) {
	if len(args) == 0 {
		e.Msg(m)
	} else {
		e.Msgf(m, args...)
	}
}
