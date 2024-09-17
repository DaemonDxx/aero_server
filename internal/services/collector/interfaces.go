package collector

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
)

type UserDAO interface {
	Find(ctx context.Context, q entity.User) ([]entity.User, error)
}

type NotificationService interface {
	ActualOrderNotify(userID uint, o entity.Order)
	ErrorNotify(userID uint, err error)
}

type OrderService interface {
	GetActualOrder(ctx context.Context, userID uint) (entity.Order, error)
}
