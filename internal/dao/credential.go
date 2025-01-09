package dao

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var ErrCredentialNotFound = errors.New("credential not found")

type CredentialDAO struct {
	db *gorm.DB
}

func (c *CredentialDAO) FindByLogin(ctx context.Context, accLogin string, lksLogin string) (*entity.Credential, error) {
	cr := &entity.Credential{}
	if err := c.db.WithContext(ctx).First(&cr, &entity.Credential{
		AccordLogin: accLogin,
		LKSLogin:    lksLogin,
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCredentialNotFound
		} else {
			return nil, err
		}
	}
	return cr, nil
}

func (c *CredentialDAO) Save(ctx context.Context, cr *entity.Credential) error {
	return c.db.WithContext(ctx).Save(&cr).Error
}

func (c *CredentialDAO) Create(ctx context.Context, cr *entity.Credential) error {
	return c.db.WithContext(ctx).Create(&cr).Error
}

func (c *CredentialDAO) GetByID(ctx context.Context, id uint) (*entity.Credential, error) {
	cr := &entity.Credential{}
	if err := c.db.WithContext(ctx).First(cr, &entity.Credential{
		Model: gorm.Model{
			ID: id,
		},
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCredentialNotFound
		} else {
			return nil, err
		}
	}
	return cr, nil
}

func NewCredentialDAO(db *gorm.DB) *CredentialDAO {
	return &CredentialDAO{
		db: db,
	}
}
