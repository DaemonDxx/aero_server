package logger

import (
	"github.com/rs/zerolog"
	"os"
)

type Env string

var (
	DEV  Env = "DEV"
	PROD Env = "PROD"
)

func NewLogger(e Env) *zerolog.Logger {
	var l zerolog.Logger
	if e == DEV {
		l = zerolog.New(zerolog.ConsoleWriter{
			Out: os.Stdout,
		}).
			Level(zerolog.DebugLevel).
			With().
			Timestamp().
			Str("app", "backend").
			Logger()
		return &l
	} else if e == PROD {
		l = zerolog.New(os.Stdout)
		return &l
	}
	l = zerolog.New(os.Stdout)
	return &l
}
