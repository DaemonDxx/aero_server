package puller

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/services"
	"github.com/rs/zerolog"
)

const servName = "collector "

type OrderDAO interface {
	Save(ctx context.Context, o *entity.Order) error
	FindLastOrders(ctx context.Context, cr *entity.Credential, limit int) ([]entity.Order, error)
}

type LKSApi interface {
	GetActualDuty(ctx context.Context, cr *entity.Credential) ([]entity.OrderItem, error)
}

type Service struct {
	services.LoggedService
	dao OrderDAO
	api LKSApi
}

func NewPullerService(d OrderDAO, a LKSApi, log *zerolog.Logger) *Service {
	return &Service{
		LoggedService: services.NewLoggedService("puller_service", log),
		dao:           d,
		api:           a,
	}
}
