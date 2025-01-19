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
	CredentialID uint        `json:"credentialID" gorm:"column:credential_id"`
	Items        []OrderItem `json:"items"`
	Status       OrderStatus `json:"status"`
}
