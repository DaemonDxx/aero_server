package dao

import (
	"fmt"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/config"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	db  *gorm.DB
	log *zerolog.Logger
}

func NewDatabase(cfg config.DBConfig) (*Database, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect to db error: %e", err)
	}
	return &Database{db: db}, err
}

func (d *Database) AutoMigrate() error {
	if err := d.db.AutoMigrate(&entity.Account{}); err != nil {
		return fmt.Errorf("migrate account entity error: %e", err)
	}
	return nil
}

func (d *Database) GetInstance() *gorm.DB {
	return d.db
}
