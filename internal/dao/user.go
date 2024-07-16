package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/daemondxx/lks_back/entity"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Create(ctx context.Context, user *entity.User) error {
	return dao.db.WithContext(ctx).Create(&user).Error
}

func (dao *UserDAO) GetByID(ctx context.Context, id uint) (entity.User, error) {
	var user entity.User
	if err := dao.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, ErrUserNotFound
		} else {
			return user, nil
		}
	}
	return user, nil
}

func (dao *UserDAO) Find(ctx context.Context, q *entity.User) ([]entity.User, error) {
	var users []entity.User
	if err := dao.db.WithContext(ctx).Where(q).Find(&users).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find user error %e", err)
		}
	}
	return users, nil
}

func (dao *UserDAO) Update(ctx context.Context, user entity.User) error {
	return dao.db.WithContext(ctx).Model(&user).Updates(&entity.User{
		AccordLogin:    user.AccordLogin,
		LKSLogin:       user.LKSLogin,
		AccordPassword: user.AccordPassword,
		LKSPassword:    user.LKSPassword,
		IsActive:       user.IsActive,
	}).Error
}
