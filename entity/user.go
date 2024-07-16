package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID             uint
	AccordLogin    string `gorm:"unique"`
	LKSLogin       string `gorm:"unique"`
	AccordPassword string
	LKSPassword    string
	IsActive       bool
}
