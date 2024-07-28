package lks

import "time"

type archiveDutyResponse struct {
	Status string `json:"status"`
	Model  struct {
		StaffNumber    string        `json:"StaffNumber"`
		DateFrom       time.Time     `json:"DateFrom"`
		DateTo         time.Time     `json:"DateTo"`
		TotalPages     int           `json:"TotalPages"`
		AchievedDuties []ArchiveDuty `json:"AchievedDuties"`
	} `json:"model"`
	Timestamp time.Time `json:"timestamp"`
}

type ArchiveDuty struct {
	StartTime        time.Time `json:"StartDateTime"`
	EndTime          time.Time `json:"EndDateTime"`
	FlightNumber     string    `json:"FlightNumber"`
	Aircraft         string    `json:"AircraftType"`
	FlightDuration   int       `json:"FlightDuration"`
	WorkTimeDuration int       `json:"WorkTimeDuration"`
	IsPass           string    `json:"IsPass"`
}
