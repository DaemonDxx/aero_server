package entity

import (
	"gorm.io/gorm"
)

type OrderStatus int

const (
	AwaitConfirmation OrderStatus = iota
	Confirm
	Limited
)

type Order struct {
	gorm.Model
	UserID uint
	Items  []OrderItem
	Status OrderStatus
}
