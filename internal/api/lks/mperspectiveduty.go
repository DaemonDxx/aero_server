package lks

import "time"

type monthPerspectiveResponse struct {
	Status string `json:"status"`
	Model  struct {
		PerspectivePlan []PerspectiveDuty `json:"PerspectivePlan"`
		ConfDate        time.Time         `json:"ConfDate"`
		CanBeConfirmed  bool              `json:"CanBeConfirmed"`
	} `json:"model"`
}

type PerspectiveDuty struct {
	StartDate    time.Time `json:"DateBeg"`
	EndDate      time.Time `json:"DateEnd"`
	FlightNumber string    `json:"NFlight"`
	Route        string    `json:"Marsh"`
}
