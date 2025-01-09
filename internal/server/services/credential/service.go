package service_credential

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/services"
	"github.com/rs/zerolog"
)

const servName = "credential_service"

type AccountService interface {
	ConnectCredential(ctx context.Context, acc *entity.Account, cred *entity.Credential) error
}

type AuthChecker interface {
	Check(ctx context.Context, accLogin string, accPass string, lksLogin string, lksPass string) error
}

type CredentialDAO interface {
	FindByLogin(ctx context.Context, accLogin string, lksLogin string) (*entity.Credential, error)
	Save(ctx context.Context, cred *entity.Credential) error
	Create(ctx context.Context, cred *entity.Credential) error
	GetByID(ctx context.Context, id uint) (*entity.Credential, error)
}

type Service struct {
	services.LoggedService
	checker AuthChecker
	dao     CredentialDAO
	accServ AccountService
}

func NewCredentialService(accServ AccountService, auth AuthChecker, dao CredentialDAO, log *zerolog.Logger) *Service {
	return &Service{
		LoggedService: services.NewLoggedService("credential_service", log),
		checker:       auth,
		dao:           dao,
		accServ:       accServ,
	}
}
