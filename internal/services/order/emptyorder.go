package order

import (
	"github.com/daemondxx/lks_back/entity"
)

func newEmptyOrder(userID uint) entity.Order {
	return entity.Order{
		UserID: userID,
		Items:  make([]entity.OrderItem, 0),
		Status: entity.Confirm,
	}
}
