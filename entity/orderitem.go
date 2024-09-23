package entity

import (
	"gorm.io/gorm"
	"time"
)

type OrderItem struct {
	gorm.Model  `json:"gorm.Model"`
	Flights     []Flight   `gorm:"foreignKey:OItemID" json:"flights,omitempty"`
	Departure   time.Time  `json:"departure"`
	Arrival     time.Time  `json:"arrival"`
	Description string     `json:"description,omitempty"`
	Route       string     `json:"route,omitempty"`
	ConfirmDate *time.Time `json:"confirmDate,omitempty"`
	OrderID     uint       `json:"orderID,omitempty"`
}
