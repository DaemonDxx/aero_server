package service_account

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/services"
)

func (s *Service) GetByID(ctx context.Context, id uint) (*entity.Account, error) {
	l := s.GetLogger("get_by_id")

	acc, err := s.dao.Get(ctx, id)
	if err != nil {
		l.Error().Err(err).Msg("find account by id error")
		return nil, &services.ErrServ{
			Service: servName,
			Message: "get account by id error",
			Err:     err,
		}
	}
	return acc, nil
}
