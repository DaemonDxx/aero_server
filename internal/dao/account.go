package dao

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var ErrAccountNotFound = errors.New("account not found")

type AccountDAO struct {
	db *gorm.DB
}

func NewAccountDAO(db *gorm.DB) *AccountDAO {
	return &AccountDAO{
		db: db,
	}
}

func (a *AccountDAO) Find(ctx context.Context, tgID uint64) ([]entity.Account, error) {
	var accs []entity.Account
	if err := a.db.WithContext(ctx).Find(&accs, entity.Account{TelegramID: tgID}).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return accs, err
	}
	return accs, nil
}

func (a *AccountDAO) Create(ctx context.Context, acc *entity.Account) error {
	return a.db.WithContext(ctx).Create(&acc).Error
}

func (a *AccountDAO) Save(ctx context.Context, acc *entity.Account) error {
	return a.db.WithContext(ctx).Save(&acc).Error
}

func (a *AccountDAO) Get(ctx context.Context, id uint) (*entity.Account, error) {
	var acc *entity.Account
	if err := a.db.WithContext(ctx).First(acc, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAccountNotFound
		} else {
			return nil, err
		}
	}
	return acc, nil
}
