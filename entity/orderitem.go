package entity

import (
	"gorm.io/gorm"
	"time"
)

type OrderItem struct {
	gorm.Model
	Flights     []Flight `gorm:"foreignKey:OItemID"`
	Departure   time.Time
	Arrival     time.Time
	Description string
	Route       string
	ConfirmDate *time.Time
	OrderID     uint
}
