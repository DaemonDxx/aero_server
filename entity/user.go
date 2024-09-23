package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID             uint   `json:"id,omitempty"`
	AccordLogin    string `gorm:"unique" json:"accordLogin,omitempty"`
	LKSLogin       string `gorm:"unique" json:"lksLogin,omitempty"`
	AccordPassword string `json:"accordPassword,omitempty"`
	LKSPassword    string `json:"lksPassword,omitempty"`
	IsActive       bool   `json:"isActive,omitempty"`
}
