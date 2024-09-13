package servers

import (
	"context"
	"errors"
	authpb "github.com/daemondxx/lks_back/gen/pb/go/auth"
	userpb "github.com/daemondxx/lks_back/gen/pb/go/user"
	"github.com/daemondxx/lks_back/internal/api/lks"
	"github.com/daemondxx/lks_back/internal/services/user"
	"github.com/rs/zerolog"
	"os"
)

type AuthChecker interface {
	Check(ctx context.Context, accLogin string, accPass string, lksLogin string, lksPass string) error
}

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	c   AuthChecker
	u   UserService
	log *zerolog.Logger
}

func NewAuthServer(u UserService, c AuthChecker, log *zerolog.Logger) *AuthServer {
	if log == nil {
		var l zerolog.Logger
		l = zerolog.New(os.Stdout).Level(zerolog.NoLevel)
		log = &l
	}

	return &AuthServer{
		c:   c,
		u:   u,
		log: log,
	}
}

func (a *AuthServer) Auth(ctx context.Context, r *userpb.UserInfo) (*authpb.AuthResponse, error) {
	u, err := a.u.Register(ctx, r.AccordLogin, r.AccordPassword, r.LksLogin, r.LksPassword)
	if err != nil {
		if errors.Is(err, user.ErrUserIsRegister) {
			u, err = a.u.GetUserByAccordLogin(ctx, r.AccordLogin)
			if err != nil {
				return nil, ErrInternal
			}
			return &authpb.AuthResponse{
				User: &userpb.User{
					Id:       uint64(u.ID),
					IsActive: u.IsActive,
				}}, nil
		} else {
			return nil, ErrInternal
		}
	}
	return &authpb.AuthResponse{User: &userpb.User{
		Id:       uint64(u.ID),
		IsActive: u.IsActive,
	}}, nil
}

func (a *AuthServer) Check(ctx context.Context, r *userpb.UserInfo) (*authpb.CheckResponse, error) {
	log := a.getLogger("check")
	log.Debug().Msg("start check auth info")

	if err := a.c.Check(ctx, r.AccordLogin, r.AccordPassword, r.LksLogin, r.LksPassword); err != nil {
		if !(errors.Is(err, lks.ErrAccordAuth) || errors.Is(err, lks.ErrLKSAuth)) {
			return nil, ErrInternal
		} else {
			var system authpb.AuthSystem
			if errors.Is(err, lks.ErrAccordAuth) {
				system = authpb.AuthSystem_SYSTEM_ACCORD
			} else {
				system = authpb.AuthSystem_SYSTEM_LKS
			}
			return &authpb.CheckResponse{Details: &authpb.ErrorDetails{System: system}}, nil
		}
	}
	return &authpb.CheckResponse{Details: &authpb.ErrorDetails{System: authpb.AuthSystem_SYSTEM_UNKNOWN}}, nil
}

func (a *AuthServer) getLogger(method string) zerolog.Logger {
	return a.log.With().Str("method", method).Logger()
}

func (a *AuthServer) mustEmbedUnimplementedAuthServiceServer() {
	panic("implement me")
}
