package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/services"
	"github.com/rs/zerolog"
)

const servName = "userService"

var (
	ErrUserIsRegister = errors.New("user is register in system")
	ErrUserNotFound   = errors.New("user not found")
)

type DAO interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uint) (entity.User, error)
	Find(ctx context.Context, q *entity.User) ([]entity.User, error)
	Update(ctx context.Context, user entity.User) error
}

type Service struct {
	d   DAO
	log *zerolog.Logger
}

func NewUserService(d DAO, log *zerolog.Logger) *Service {
	l := log.With().Str("service", "user_service").Logger()
	return &Service{
		d:   d,
		log: &l,
	}
}

func (s *Service) Register(ctx context.Context, accLogin string, accPass string, lksLogin string, lksPass string) (entity.User, error) {
	var user entity.User
	log := s.getLogger("register")

	users, err := s.d.Find(ctx, &entity.User{
		AccordLogin: accLogin,
	})

	if err != nil {
		log.
			Debug().
			Msgf("find user by accord login error: %e", err)

		return user, &services.ErrServ{
			Service: servName,
			Message: "find user with accord login error",
			Err:     err,
		}
	}

	if len(users) != 0 {
		log.
			Info().
			Msgf("user with accord login %s is registred", accLogin)
		return user, ErrUserIsRegister
	}

	user.AccordLogin = accLogin
	user.LKSLogin = lksLogin
	user.AccordPassword = accPass
	user.LKSPassword = lksPass
	user.IsActive = true

	if err := s.d.Create(ctx, &user); err != nil {
		log.
			Debug().
			Msgf("create new user error: %e", err)

		return user, &services.ErrServ{
			Service: servName,
			Message: "create user error",
			Err:     err,
		}
	}

	return user, nil
}

func (s *Service) UpdateAccord(ctx context.Context, userID uint, login string, password string) error {
	log := s.getLogger("update_accord")

	if err := s.hasUserByID(ctx, userID); err != nil {
		return err
	}

	if err := s.d.Update(ctx, entity.User{
		ID:             userID,
		AccordLogin:    login,
		AccordPassword: password,
	}); err != nil {
		log.Debug().
			Uint("user_id", userID).
			Msgf("update accord info error: %e", err)
		return &services.ErrServ{
			Service: servName,
			Message: "update accord auth info error",
			Err:     err,
		}
	}
	return nil
}

func (s *Service) UpdateLKS(ctx context.Context, userID uint, login string, password string) error {
	log := s.getLogger("update_lks")

	if err := s.hasUserByID(ctx, userID); err != nil {
		return err
	}

	if err := s.d.Update(ctx, entity.User{
		ID:             userID,
		AccordLogin:    login,
		AccordPassword: password,
	}); err != nil {
		log.Debug().
			Uint("user_id", userID).
			Msgf("update lks info error: %e", err)
		return &services.ErrServ{
			Service: servName,
			Message: "update lks info error",
			Err:     err,
		}
	}
	return nil
}

func (s *Service) UpdateActiveStatus(ctx context.Context, userID uint, status bool) error {
	log := s.getLogger("update_active_status")
	if err := s.hasUserByID(ctx, userID); err != nil {
		return err
	}

	if err := s.d.Update(ctx, entity.User{
		ID:       userID,
		IsActive: status,
	}); err != nil {
		log.Debug().
			Uint("user_id", userID).
			Msgf("update user status error: %e", err)
		return &services.ErrServ{
			Service: servName,
			Message: "update user status error",
			Err:     err,
		}
	}
	return nil
}

func (s *Service) GetUserByAccordLogin(ctx context.Context, accLogin string) (entity.User, error) {
	var u entity.User
	log := s.getLogger("get_user_by_accord_login")

	users, err := s.d.Find(ctx, &entity.User{AccordLogin: accLogin})
	if err != nil {
		log.Debug().
			Str("login", accLogin).
			Msgf("find user by accord login error: %e", err)
		return u, &services.ErrServ{
			Service: servName,
			Message: "find user by accord login error",
			Err:     err,
		}
	}

	if len(users) == 0 {
		return u, ErrUserNotFound
	}

	u = users[0]
	return u, nil
}

func (s *Service) GetUserByID(ctx context.Context, id uint) (entity.User, error) {
	log := s.getLogger("get_user_by_id")

	if u, err := s.d.GetByID(ctx, id); err != nil {
		log.Debug().
			Uint("user_id", id).
			Msgf("find user by login error: %e", err)
		return u, &services.ErrServ{
			Service: servName,
			Message: "get user by id error",
			Err:     err,
		}
	} else {
		return u, nil
	}
}

func (s *Service) hasUserByID(ctx context.Context, userID uint) error {
	log := s.getLogger("has_user_by_id")
	if _, err := s.d.GetByID(ctx, userID); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return &services.ErrServ{
				Service: servName,
				Message: fmt.Sprintf("user with id %d not found", userID),
				Err:     err,
			}
		} else {
			log.Debug().
				Uint("user_id", userID).
				Msgf("get user by id error: %e", err)
			return &services.ErrServ{
				Service: servName,
				Message: "find user by id error",
				Err:     err,
			}
		}
	}
	return nil
}

func (s *Service) getLogger(method string) zerolog.Logger {
	return s.log.With().Str("method", method).Logger()
}
