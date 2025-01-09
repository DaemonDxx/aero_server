package service_account

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/services"
)

func (s *Service) CreateAccount(ctx context.Context, tgID uint64) (string, error) {
	log := s.GetLogger("create_account")

	if ok, err := s.isAccountCreated(ctx, tgID); ok || err != nil {
		if err != nil {
			log.Error().Err(err).Msg("failed to create account")
			return "", &services.ErrServ{
				Service: servName,
				Message: "check is account created error",
				Err:     err,
			}
		} else {
			return "", ErrAccountExists
		}
	}

	acc := &entity.Account{
		TelegramID: tgID,
	}

	if err := s.dao.Create(ctx, acc); err != nil {
		log.Error().Err(err).Msg("failed to create account")
		return "", &services.ErrServ{
			Service: servName,
			Message: "failed to create account",
			Err:     err,
		}
	}

	return s.tokenServ.Create(acc), nil
}
