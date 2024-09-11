package notifier

import (
	"context"
	"encoding/json"
	"github.com/daemondxx/lks_back/internal/services"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

const servName = "notification_service"

type Service struct {
	services.LoggedService
	w *kafka.Writer
}

func NewNotifierService(w *kafka.Writer, log *zerolog.Logger) *Service {
	l := log.With().Str("service", "notifier_service").Logger()
	return &Service{
		w:             w,
		LoggedService: services.NewLoggedService(&l),
	}
}

func (s Service) Notify(ctx context.Context, n Notification) error {
	b, err := json.Marshal(n)
	if err != nil {
		return &services.ErrServ{
			Service: servName,
			Message: "marshal notification object error",
			Err:     err,
		}
	}

	if err := s.w.WriteMessages(ctx, kafka.Message{Key: []byte(n.Key), Value: b}); err != nil {
		return &services.ErrServ{
			Service: servName,
			Message: "send message error",
			Err:     err,
		}
	}

	return nil
}
