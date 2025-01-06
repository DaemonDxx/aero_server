package service_account

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/services"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const servName = "account_service"

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrAccountExists   = errors.New("account already exists")
)

type AccountDAO interface {
	Find(ctx context.Context, tgID uint64) ([]entity.Account, error)
	Create(ctx context.Context, a *entity.Account) error
	Save(ctx context.Context, a *entity.Account) error
	Get(ctx context.Context, id uint) (*entity.Account, error)
}

type TokenService interface {
	Create(acc *entity.Account) string
}

type Service struct {
	services.LoggedService
	dao       AccountDAO
	tokenServ TokenService
}

func NewAccountService(dao AccountDAO, tServ TokenService, log *zerolog.Logger) *Service {
	return &Service{
		LoggedService: services.NewLoggedService("account_service", log),
		dao:           dao,
		tokenServ:     tServ,
	}
}
