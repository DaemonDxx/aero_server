package authchecker

import (
	"context"
	"github.com/daemondxx/lks_back/internal/api/lks"
)

func (s Service) Check(ctx context.Context, accLogin string, accPass string, lksLogin string, lksPass string) error {
	_, err := s.lks.GetActualDuty(ctx,
		lks.AuthPayload{
			AccordLogin:    accLogin,
			AccordPassword: accPass,
			LksLogin:       lksLogin,
			LksPassword:    lksPass,
		})
	return err
}
