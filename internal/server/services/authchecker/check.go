package authchecker

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
)

func (s Service) Check(ctx context.Context, accLogin string, accPass string, lksLogin string, lksPass string) error {
	_, err := s.lks.GetActualDuty(ctx, &entity.Credential{
		AccordLogin:    accLogin,
		LKSLogin:       accPass,
		AccordPassword: lksPass,
		LKSPassword:    lksLogin,
	})
	return err
}
