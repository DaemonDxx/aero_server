package services

import (
	"github.com/rs/zerolog"
	"os"
)

type LoggedService struct {
	log *zerolog.Logger
}

func NewLoggedService(servName string, log *zerolog.Logger) LoggedService {
	var l zerolog.Logger

	if log == nil {
		var l zerolog.Logger
		l = zerolog.New(os.Stdout).Level(zerolog.NoLevel)
		log = &l
	} else {
		l = log.With().Str("service", servName).Logger()
	}

	return LoggedService{
		log: &l,
	}
}

func (a *LoggedService) GetLogger(method string) zerolog.Logger {
	return a.log.With().Str("method", method).Logger()
}
