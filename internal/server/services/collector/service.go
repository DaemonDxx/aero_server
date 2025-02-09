package collector

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/server/services/notifier"
	"github.com/daemondxx/lks_back/internal/services"
	"github.com/rs/zerolog"
	"time"
)

const servName = "collector_service"

type CredentialDAO interface {
	GetActualCredential(ctx context.Context) ([]entity.Credential, error)
	SetActiveStatus(ctx context.Context, cr *entity.Credential) error
	SetInactiveStatus(ctx context.Context, cr *entity.Credential) error
}

type OrderService interface {
	PullNewOrder(ctx context.Context, cr *entity.Credential) (*entity.Order, error)
	Create(ctx context.Context, cr *entity.Credential, i []entity.OrderItem) (*entity.Order, error)
}

type NotificationService interface {
	Send(msg notifier.Message) error
}

type Config struct {
	MaxAttempts     int
	MinTimeoutRetry time.Duration
}

type Service struct {
	services.LoggedService
	crDAO  CredentialDAO
	oServ  OrderService
	notify NotificationService
	cfg    Config
}

func NewCollectorService(
	dao CredentialDAO,
	oServ OrderService,
	n NotificationService,
	c Config,
	log *zerolog.Logger,
) *Service {
	if c.MaxAttempts == 0 {
		c.MaxAttempts = defaultMaxAttempts
	}
	if c.MinTimeoutRetry == 0 {
		c.MinTimeoutRetry = defaultMinTimeoutRetry
	}

	return &Service{
		LoggedService: services.NewLoggedService("collector_service", log),
		crDAO:         dao,
		oServ:         oServ,
		notify:        n,
		cfg:           c,
	}
}
