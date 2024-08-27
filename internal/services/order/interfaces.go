package order

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/api/lks"
)

type DAO interface {
	Create(ctx context.Context, o *entity.Order) error
	GetLastOrder(ctx context.Context, userID uint) (entity.Order, error)
	Save(ctx context.Context, o *entity.Order) error
}

type LksAPI interface {
	GetActualDuty(ctx context.Context, p lks.AuthPayload) ([]lks.CurrentDuty, error)
	GetPerspectiveDuty(ctx context.Context, p lks.AuthPayload, month int, year int) ([]lks.PerspectiveDuty, error)
	GetArchiveDuty(ctx context.Context, p lks.AuthPayload, month int, year int) ([]lks.ArchiveDuty, error)
}

type UserService interface {
	GetUserByID(ctx context.Context, id uint) (entity.User, error)
}
