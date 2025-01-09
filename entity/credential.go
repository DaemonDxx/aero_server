package entity

import "gorm.io/gorm"

type Credential struct {
	gorm.Model
	AccordLogin    string    `gorm:"unique" json:"accordLogin,omitempty"`
	LKSLogin       string    `gorm:"unique" json:"lksLogin,omitempty"`
	AccordPassword string    `json:"accordPassword,omitempty"`
	LKSPassword    string    `json:"lksPassword,omitempty"`
	Accounts       []Account `gorm:"foreignKey:CredentialID" json:"accounts,omitempty"`
	IsActual       bool      `json:"isActual"`
}
