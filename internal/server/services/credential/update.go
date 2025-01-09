package service_credential

import (
	"context"
	"github.com/daemondxx/lks_back/internal/services"
)

func (s *Service) Update(ctx context.Context, id uint, accPass string, lksPass string) error {
	cr, err := s.dao.GetByID(ctx, id)
	if err != nil {
		return &services.ErrServ{
			Service: servName,
			Message: "get credential by id error",
			Err:     err,
		}
	}

	if err := s.checker.Check(ctx, cr.AccordLogin, accPass, cr.LKSLogin, lksPass); err != nil {
		return &services.ErrServ{
			Service: servName,
			Message: "check credential error",
			Err:     err,
		}
	}

	cr.AccordPassword = accPass
	cr.LKSPassword = lksPass
	cr.IsActual = true

	if err := s.dao.Save(ctx, cr); err != nil {
		return &services.ErrServ{
			Service: servName,
			Message: "save credential error",
			Err:     err,
		}
	}

	return nil
}
