package service_account

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
)

func (s *Service) ConnectCredential(ctx context.Context, acc *entity.Account, cr *entity.Credential) error {
	acc.CredentialID = &cr.ID
	return s.dao.Save(ctx, acc)
}
