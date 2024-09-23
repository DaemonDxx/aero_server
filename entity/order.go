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
	UserID uint        `json:"userID"`
	Items  []OrderItem `json:"items"`
	Status OrderStatus `json:"status"`
}
