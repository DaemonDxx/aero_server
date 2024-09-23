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
	gorm.Model   `json:"gorm.Model"`
	FlightNumber string         `json:"flightNumber,omitempty"`
	Airplane     string         `json:"airplane,omitempty"`
	Departure    *time.Time     `json:"departure,omitempty"`
	Arrival      *time.Time     `json:"arrival,omitempty"`
	Duration     *time.Duration `json:"duration,omitempty"`
	Status       FlightStatus   `json:"status,omitempty"`
	OItemID      uint           `json:"oItemID,omitempty"`
}
