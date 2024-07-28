package lks

import (
	"time"
)

type currentDutyResponse struct {
	Status string `json:"status"`
	Model  struct {
		Duties         []CurrentDuty `json:"Duties"`
		Message        string        `json:"Message"`
		CanBeConfirmed bool          `json:"CanBeConfirmed"`
	} `json:"model"`
}

type CurrentDuty struct {
	Code         string      `json:"Code"`
	FlightNumber string      `json:"FlightNumber"`
	AircraftType string      `json:"AircraftType"`
	Route        string      `json:"Route"`
	StartDate    time.Time   `json:"StartDate"`
	EndDate      time.Time   `json:"EndDate"`
	ConfirmDate  time.Time   `json:"ConfirmDate"`
	ConfirmType  int         `json:"ConfirmType"`
	Target       string      `json:"Target"`
	Place        string      `json:"Place"`
	Note         string      `json:"Note"`
	BlockDate    interface{} `json:"BlockDate"`
}
