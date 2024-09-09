package servers

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
)

type UserService interface {
	Register(ctx context.Context, accLogin string, accPass string, lksLogin string, lksPass string) (entity.User, error)
	UpdateAccord(ctx context.Context, userID uint, login string, password string) error
	UpdateLKS(ctx context.Context, userID uint, login string, password string) error
	UpdateActiveStatus(ctx context.Context, userID uint, status bool) error
	GetUserByAccordLogin(ctx context.Context, accLogin string) (entity.User, error)
	GetUserByID(ctx context.Context, id uint) (entity.User, error)
}
type OrderService interface {
	GetActualOrder(ctx context.Context, userID uint) (entity.Order, error)
}
