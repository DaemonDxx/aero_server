package services

import "github.com/rs/zerolog"

type LoggedService struct {
	log *zerolog.Logger
}

func NewLoggedService(log *zerolog.Logger) LoggedService {
	return LoggedService{
		log: log,
	}
}

func (a *LoggedService) GetLogger(method string) zerolog.Logger {
	return a.log.With().Str("method", method).Logger()
}
