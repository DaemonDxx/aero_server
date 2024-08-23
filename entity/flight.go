package entity

import (
	"gorm.io/gorm"
	"time"
)

type FlightStatus int

const (
	Await FlightStatus = iota
	Completed
	Canceled
)

type Flight struct {
	gorm.Model
	FlightNumber string
	Airplane     string
	Departure    *time.Time
	Arrival      *time.Time
	Duration     *time.Duration
	Status       FlightStatus
	OItemID      uint
}
