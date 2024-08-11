package entity

import (
	"time"
)

type FlightInfo struct {
	FlightNumber  string `gorm:"primaryKey;unique"`
	From          string
	To            string
	TimeDeparture time.Time
	Duration      time.Duration
	AvgDuration   time.Duration
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}
