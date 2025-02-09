package lks_mock

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	lks "github.com/daemondxx/lks_back/internal/api/lks"
)

type ActualFn func(ctx context.Context, cr *entity.Credential) ([]entity.OrderItem, error)
type PerspectiveFn func(ctx context.Context, cr *entity.Credential, month int, year int) ([]lks.PerspectiveDuty, error)
type ArchiveFn func(ctx context.Context, cr *entity.Credential, month int, year int) ([]lks.ArchiveDuty, error)

type APIMock struct {
	actualFn      ActualFn
	perspectiveFn PerspectiveFn
	archiveFn     ArchiveFn
}

func NewLKSApiMock(actualFn ActualFn) *APIMock {
	return &APIMock{
		actualFn: actualFn,
	}
}

func (a *APIMock) GetActualDuty(ctx context.Context, cr *entity.Credential) ([]entity.OrderItem, error) {
	return a.actualFn(ctx, cr)
}

func (a *APIMock) GetPerspectiveDuty(ctx context.Context, cr *entity.Credential, month int, year int) ([]lks.PerspectiveDuty, error) {
	return a.perspectiveFn(ctx, cr, month, year)
}

func (a *APIMock) GetArchiveDuty(ctx context.Context, cr *entity.Credential, month int, year int) ([]lks.ArchiveDuty, error) {
	return a.archiveFn(ctx, cr, month, year)
}
