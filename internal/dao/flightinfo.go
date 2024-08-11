package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/daemondxx/lks_back/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrNotFoundInfo = errors.New("flight info not found")
)

type FlightInfoDAO struct {
	db *gorm.DB
}

func NewFlightInfoDAO(db *gorm.DB) *FlightInfoDAO {
	return &FlightInfoDAO{
		db: db,
	}
}

func (f *FlightInfoDAO) GetByNumberFlight(ctx context.Context, number uint) (entity.FlightInfo, error) {
	//todo если данные старше недели - удалить
	var info entity.FlightInfo
	if err := f.db.WithContext(ctx).First(&info, number).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return info, ErrNotFoundInfo
		} else {
			return info, fmt.Errorf("get flight info error: %w", err)
		}
	}
	return info, nil
}

func (f *FlightInfoDAO) Save(ctx context.Context, info *entity.FlightInfo) error {
	if err := f.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "flight_number"}},
		UpdateAll: true,
	}).Create(info).Error; err != nil {
		return fmt.Errorf("save flight info error: %w", err)
	}
	return nil
}
