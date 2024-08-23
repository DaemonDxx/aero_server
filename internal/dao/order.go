package dao

import (
	"context"
	"fmt"
	"github.com/daemondxx/lks_back/entity"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var ErrOrderNotFound = errors.New("order not found")

type OrderDAO struct {
	db *gorm.DB
}

func NewOrderDAO(db *gorm.DB) *OrderDAO {
	return &OrderDAO{
		db: db,
	}
}

func (d *OrderDAO) Create(ctx context.Context, o *entity.Order) error {
	tx := d.db.Begin()
	if err := tx.WithContext(ctx).Create(o).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("create order error: %e", err)
	}
	tx.Commit()
	return nil
}

func (d *OrderDAO) GetLastOrder(ctx context.Context, userID uint) (entity.Order, error) {
	var order entity.Order
	if err := d.db.
		WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(1).
		Take(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return order, ErrOrderNotFound
		} else {
			return order, err
		}
	}
	if err := d.db.
		WithContext(ctx).
		Model(&order).
		Association("Items").
		Find(&order.Items); err != nil {
		return order, err
	}

	for i, item := range order.Items {
		if err := d.db.
			WithContext(ctx).
			Model(&item).
			Association("Flights").
			Find(&order.Items[i].Flights); err != nil {
			return order, fmt.Errorf("get flight by item %d error: %e", item.ID, err)
		}
	}

	return order, nil
}

func (d *OrderDAO) Save(ctx context.Context, order *entity.Order) error {
	return d.db.
		WithContext(ctx).
		Save(order).
		Error
}
