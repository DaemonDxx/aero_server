package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/services"
	"gorm.io/gorm"
)

const serName = "userService"

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
	d DAO
}

func NewUserService(d DAO) *Service {
	return &Service{
		d: d,
	}
}

func (s *Service) Register(ctx context.Context, accLogin string, accPass string, lksLogin string, lksPass string) (entity.User, error) {
	var user entity.User

	users, err := s.d.Find(ctx, &entity.User{
		AccordLogin: accLogin,
	})

	if err != nil {
		return user, &services.ErrServ{
			Service: serName,
			Message: "find user with accord login error",
			Err:     err,
		}
	}

	if len(users) != 0 {
		return user, ErrUserIsRegister
	}

	user.AccordLogin = accLogin
	user.LKSLogin = lksLogin
	user.AccordPassword = accPass
	user.LKSPassword = lksPass
	user.IsActive = true

	if err := s.d.Create(ctx, &user); err != nil {
		return user, &services.ErrServ{
			Service: serName,
			Message: "create user error",
			Err:     err,
		}
	}

	return user, nil
}

func (s *Service) UpdateAccord(ctx context.Context, userID uint, login string, password string) error {
	if err := s.d.Update(ctx, entity.User{
		ID:             userID,
		AccordLogin:    login,
		AccordPassword: password,
	}); err != nil {
		return &services.ErrServ{
			Service: serName,
			Message: "update accord auth info error",
			Err:     err,
		}
	}
	return nil
}

func (s *Service) UpdateLKS(ctx context.Context, userID uint, login string, password string) error {
	if err := s.d.Update(ctx, entity.User{
		ID:             userID,
		AccordLogin:    login,
		AccordPassword: password,
	}); err != nil {
		return fmt.Errorf("update accord auth info error: %e", err)
	}
	return nil
}

func (s *Service) UpdateActiveStatus(ctx context.Context, userID uint, status bool) error {
	if err := s.d.Update(ctx, entity.User{
		ID:       userID,
		IsActive: status,
	}); err != nil {
		return fmt.Errorf("update user status error: %w", err)
	}
	return nil
}

func (s *Service) GetUserByAccordLogin(ctx context.Context, accLogin string) (entity.User, error) {
	user, err := s.d.Find(ctx, &entity.User{AccordLogin: accLogin})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, ErrUserNotFound
		}
		return entity.User{}, fmt.Errorf("find user by id error: %e", err)
	}
	return user[0], nil
}
