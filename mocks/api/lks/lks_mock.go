package lks_mock

import (
	"context"
	lks "github.com/daemondxx/lks_back/internal/api/lks"
)

type ActualFn func(ctx context.Context, p lks.AuthPayload) ([]lks.CurrentDuty, error)
type PerspectiveFn func(ctx context.Context, p lks.AuthPayload, month int, year int) ([]lks.PerspectiveDuty, error)
type ArchiveFn func(ctx context.Context, p lks.AuthPayload, month int, year int) ([]lks.ArchiveDuty, error)

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

func (a *APIMock) GetActualDuty(ctx context.Context, p lks.AuthPayload) ([]lks.CurrentDuty, error) {
	return a.actualFn(ctx, p)
}

func (a *APIMock) GetPerspectiveDuty(ctx context.Context, p lks.AuthPayload, month int, year int) ([]lks.PerspectiveDuty, error) {
	return a.perspectiveFn(ctx, p, month, year)
}

func (a *APIMock) GetArchiveDuty(ctx context.Context, p lks.AuthPayload, month int, year int) ([]lks.ArchiveDuty, error) {
	return a.archiveFn(ctx, p, month, year)
}
