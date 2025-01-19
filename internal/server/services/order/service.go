package order

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/services"
	"github.com/rs/zerolog"
)

const servName = "order_service"

type OrderDAO interface {
	FindLastOrders(ctx context.Context, credID uint, limit int) ([]entity.Order, error)
}

type Service struct {
	services.LoggedService
	dao OrderDAO
}

func NewOrderService(dao OrderDAO, log *zerolog.Logger) *Service {
	return &Service{
		LoggedService: services.NewLoggedService(servName, log),
		dao:           dao,
	}
}
