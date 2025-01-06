package entity

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	TelegramID   uint64 `json:"telegramId"`
	CredentialID *uint
}
