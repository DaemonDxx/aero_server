package notifier

import (
	"context"
	"encoding/json"
	"github.com/daemondxx/lks_back/internal/services"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
	"time"
)

const defaultWriteTimeout = 30 * time.Second

type Service struct {
	services.LoggedService
	w *kafka.Writer
}

func NewNotifierService(w *kafka.Writer, log *zerolog.Logger) *Service {
	return &Service{
		w:             w,
		LoggedService: services.NewLoggedService(log),
	}
}

func (s Service) Notify(n Notification) {
	log := s.GetLogger("notify")

	ctx, cancel := context.WithTimeout(context.Background(), defaultWriteTimeout)
	defer cancel()

	b, err := json.Marshal(n)
	if err != nil {
		log.Err(err).Msg("marshal payload error")
	}

	if err := s.w.WriteMessages(ctx, kafka.Message{Key: []byte(n.Key), Value: b}); err != nil {
		log.Err(err).Msg("send message error")
	}
}
