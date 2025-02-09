package authchecker

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
)

type LKSApi interface {
	GetActualDuty(ctx context.Context, cr *entity.Credential) ([]entity.OrderItem, error)
}

type Service struct {
	lks LKSApi
}

func NewAuthCheckerService(lks LKSApi) *Service {
	return &Service{
		lks: lks,
	}
}
