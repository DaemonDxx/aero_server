package service_order

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/services"
	"github.com/rs/zerolog"
)

const servName = "order_service"

type OrderDAO interface {
	Save(ctx context.Context, o *entity.Order) error
	FindLastOrders(ctx context.Context, cr *entity.Credential, limit int) ([]entity.Order, error)
}

type CredentialDAO interface {
	GetByID(ctx context.Context, id uint) (*entity.Credential, error)
}

type LKSApi interface {
	GetActualDuty(ctx context.Context, cr *entity.Credential) ([]entity.OrderItem, error)
}

type Service struct {
	services.LoggedService
	orderDAO OrderDAO
	crDAO    CredentialDAO
	api      LKSApi
}

func NewOrderService(oDao OrderDAO, crDAO CredentialDAO, api LKSApi, log *zerolog.Logger) *Service {
	return &Service{
		LoggedService: services.NewLoggedService(servName, log),
		orderDAO:      oDao,
		crDAO:         crDAO,
		api:           api,
	}
}
