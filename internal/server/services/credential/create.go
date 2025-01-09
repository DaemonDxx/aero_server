package service_credential

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/api/lks"
	"github.com/daemondxx/lks_back/internal/dao"
	"github.com/daemondxx/lks_back/internal/services"
	"github.com/pkg/errors"
)

var ErrCredentialIsConnected = errors.New("credential is exist")

func (s *Service) Create(ctx context.Context, acc *entity.Account, accLogin string, accPass string, lksLogin string, lksPass string) (*entity.Credential, error) {
	l := s.GetLogger("create").With().Str("accord", accLogin).Str("lks", lksLogin).Logger()

	if acc.CredentialID != nil {
		return nil, ErrCredentialIsConnected
	}

	if err := s.checker.Check(ctx, accLogin, accPass, lksLogin, lksPass); err != nil {
		if !(errors.Is(err, lks.ErrAccordAuth) || errors.Is(err, lks.ErrLKSAuth)) {
			l.Error().Err(err).Msg("auth check error")
		}
		return nil, &services.ErrServ{
			Service: servName,
			Message: "auth check error",
			Err:     err,
		}
	}

	cr, err := s.dao.FindByLogin(ctx, accLogin, lksLogin)
	if err != nil && !errors.Is(err, dao.ErrCredentialNotFound) {
		l.Error().Err(err).Msg("find credential by login error")
		return nil, &services.ErrServ{
			Service: servName,
			Message: "find credential by login error",
			Err:     nil,
		}
	}

	if cr != nil && (accPass != cr.AccordPassword || lksPass != cr.LKSPassword) {
		cr.AccordPassword = accPass
		cr.LKSPassword = lksPass
		cr.IsActual = true

		if err := s.dao.Save(ctx, cr); err != nil {
			l.Error().Err(err).Msg("save credential error")
			return nil, &services.ErrServ{
				Service: servName,
				Message: "save credential error",
				Err:     err,
			}
		}
	} else if cr == nil {
		cr = &entity.Credential{
			AccordLogin:    accLogin,
			AccordPassword: accPass,
			LKSLogin:       lksLogin,
			LKSPassword:    lksPass,
			IsActual:       true,
		}
		if err := s.dao.Create(ctx, cr); err != nil {
			l.Error().Err(err).Msg("create credential error")
			return nil, &services.ErrServ{
				Service: servName,
				Message: "create credential error",
				Err:     err,
			}
		}
	}

	if err := s.accServ.ConnectCredential(ctx, acc, cr); err != nil {
		l.Error().Err(err).Msg("connect credential error")
		return nil, &services.ErrServ{
			Service: servName,
			Message: "connect credential error",
			Err:     err,
		}
	}

	return cr, nil
}
