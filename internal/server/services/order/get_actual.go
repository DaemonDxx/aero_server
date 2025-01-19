package order

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/services"
	"time"
)

func (s *Service) GetActualOrders(ctx context.Context, acc entity.Account) ([]entity.Order, error) {
	o, err := s.dao.FindLastOrders(ctx, *acc.CredentialID, 2)
	if err != nil {
		return nil, &services.ErrServ{
			Service: servName,
			Message: "find last order error",
			Err:     err,
		}
	}

	if len(o) != 2 {
		return o, nil
	}

	if !s.isActual(o[1]) {
		return o[:1], nil
	}

	return o, err
}

func (s *Service) isActual(o entity.Order) bool {
	if len(o.Items) == 0 {
		return false
	}

	cTime := time.Now().Unix()
	for _, i := range o.Items {
		if i.Departure.Unix() > cTime {
			return true
		}
	}

	return false
}
