package authchecker

import (
	"context"
	"github.com/daemondxx/lks_back/internal/api/lks"
)

type LKSApi interface {
	GetActualDuty(ctx context.Context, p lks.AuthPayload) ([]lks.CurrentDuty, error)
}

type Service struct {
	lks LKSApi
}

func NewAuthCheckerService(lks LKSApi) *Service {
	return &Service{
		lks: lks,
	}
}
