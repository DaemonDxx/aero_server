package service_order

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/services"
)

func (s *Service) Create(ctx context.Context, cr *entity.Credential, i []entity.OrderItem) (*entity.Order, error) {
	o := &entity.Order{
		CredentialID: cr.ID,
		Items:        i,
		Status:       entity.AwaitConfirmation,
	}
	if err := s.orderDAO.Save(ctx, o); err != nil {
		return nil, &services.ErrServ{
			Service: servName,
			Message: "save order failed",
			Err:     err,
		}
	}
	return o, nil
}
