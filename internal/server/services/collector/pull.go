package collector

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/services"
	"github.com/pkg/errors"
)

var errHasNotNewItems = errors.New("has not new items")
var ErrOrderIsExist = errors.New("order is exist")

func (s *Service) Pull(ctx context.Context, cr *entity.Credential) (*entity.Order, error) {
	i, err := s.pull(ctx, cr)
	if err != nil {
		return nil, &services.ErrServ{
			Service: servName,
			Message: "pull new items error",
			Err:     err,
		}
	}

	if len(i) == 0 {
		return s.createEmptyOrder(ctx, cr)
	}

	lo, err := s.dao.FindLastOrders(ctx, cr, 1)
	if err != nil {
		return nil, &services.ErrServ{
			Service: servName,
			Message: "find last orders error",
			Err:     err,
		}
	}

	if len(lo) == 0 || len(lo[0].Items) == 0 {
		return s.createOrder(ctx, cr, i)
	}

	i, err = s.diff(i, lo[0].Items)
	if err != nil && errors.Is(err, errHasNotNewItems) {
		return nil, ErrOrderIsExist
	}

	return s.createOrder(ctx, cr, i)
}

func (s *Service) pull(ctx context.Context, cr *entity.Credential) ([]entity.OrderItem, error) {
	i, err := s.api.GetActualDuty(ctx, cr)
	return i, err
}

// sub - last
// target - items
func (s *Service) diff(target []entity.OrderItem, sub []entity.OrderItem) ([]entity.OrderItem, error) {
	if len(target) == 0 {
		return nil, errHasNotNewItems
	}

	//если в последнем наряде нет полетов - в новый наряд берем все
	if len(sub) == 0 {
		res := make([]entity.OrderItem, len(target))
		copy(res, target)
		return res, nil
	}

	//берем последний актуальный item из наряда и если
	lItem := sub[len(sub)-1]
	if target[len(target)-1].Departure.Unix() <= lItem.Departure.Unix() {
		return nil, errHasNotNewItems
	}

	var i int
	for i = len(target) - 1; i >= 0; i-- {
		if target[i].Departure.Unix() <= lItem.Departure.Unix() {
			break
		}
	}

	target = target[i+1:]
	res := make([]entity.OrderItem, len(target))
	copy(res, target)
	return res, nil
}

func (s *Service) createEmptyOrder(ctx context.Context, cr *entity.Credential) (*entity.Order, error) {
	return s.createOrder(ctx, cr, nil)
}

func (s *Service) createOrder(ctx context.Context, cr *entity.Credential, i []entity.OrderItem) (*entity.Order, error) {
	o := &entity.Order{
		CredentialID: cr.ID,
		Items:        i,
		Status:       entity.AwaitConfirmation,
	}
	if err := s.dao.Save(ctx, o); err != nil {
		return nil, &services.ErrServ{
			Service: servName,
			Message: "save order failed",
			Err:     err,
		}
	}
	return o, nil
}
