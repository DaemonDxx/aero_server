package entity

import (
	"time"
)

type FlightInfo struct {
	FlightNumber  string        `gorm:"primaryKey;unique" json:"flightNumber,omitempty"`
	From          string        `json:"from,omitempty"`
	To            string        `json:"to,omitempty"`
	TimeDeparture time.Time     `json:"timeDeparture"`
	Duration      time.Duration `json:"duration,omitempty"`
	AvgDuration   time.Duration `json:"avgDuration,omitempty"`
	UpdatedAt     time.Time     `gorm:"autoUpdateTime" json:"updatedAt"`
}
