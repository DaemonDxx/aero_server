package lks_mock

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/api/lks"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func GetActualOrderFn_AccordLoginEquaslZero() ActualFn {
	return func(ctx context.Context, cr *entity.Credential) ([]entity.OrderItem, error) {
		n, err := strconv.ParseInt(cr.AccordLogin, 10, 64)

		if err != nil {
			return []entity.OrderItem{
				{
					Model:       gorm.Model{},
					Flights:     nil,
					Departure:   time.Date(2024, 8, 8, 16, 15, 0, 0, &time.Location{}),
					Arrival:     time.Date(2024, 8, 11, 14, 25, 0, 0, &time.Location{}),
					Description: "Явка в терминал B, 2 этаж, комната брифинга  2.12.007",
					Route:       "Ш (B) Сочи Омск Сочи Ш",
					ConfirmDate: nil,
					OrderID:     0,
				},
			}, nil
		}

		if n > 0 {
			return nil, lks.ErrAccordAuth
		} else {
			return nil, lks.ErrLKSAuth
		}
	}
}
