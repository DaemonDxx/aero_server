package lks

import (
	"fmt"
	"github.com/daemondxx/lks_back/entity"
	"github.com/rs/zerolog/log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var flightRexExp, _ = regexp.Compile("(\\d+)")

const companyFlightPrefix = "AFL"

type currentDutyResponse struct {
	Status string `json:"status"`
	Model  struct {
		Duties         []CurrentDuty `json:"Duties"`
		Message        string        `json:"Message"`
		CanBeConfirmed bool          `json:"CanBeConfirmed"`
	} `json:"model"`
}

type CurrentDuty struct {
	Code         string     `json:"Code"`
	FlightNumber string     `json:"FlightNumber"`
	AircraftType string     `json:"AircraftType"`
	Route        string     `json:"Route"`
	StartDate    time.Time  `json:"StartDate"`
	EndDate      time.Time  `json:"EndDate"`
	ConfirmDate  *time.Time `json:"ConfirmDate"`
	ConfirmType  int        `json:"ConfirmType"`
	Target       string     `json:"Target"`
	Place        string     `json:"Place"`
	Note         string     `json:"Note"`
	BlockDate    *time.Time `json:"BlockDate"`
}

func (c *currentDutyResponse) extractOrderItems() []entity.OrderItem {
	o := make([]entity.OrderItem, 0, len(c.Model.Duties))
	for _, d := range c.Model.Duties {
		item := entity.OrderItem{
			Flights:     nil,
			Departure:   d.StartDate,
			Arrival:     d.EndDate,
			Description: d.Note,
			Route:       d.Route,
			ConfirmDate: d.ConfirmDate,
		}

		flights := flightRexExp.FindAllStringSubmatch(d.FlightNumber, -1)
		if len(flights) == 0 {
			var f entity.Flight
			if strings.Contains(d.Note, "дневной") {
				f.FlightNumber = entity.DayReserve
			} else if strings.Contains(d.Note, "ночной") {
				f.FlightNumber = entity.NightReserve
			} else if strings.Contains(d.Note, "Отпуск") {
				f.FlightNumber = entity.Holyday
			} else if strings.Contains(d.Note, "Явка") {
				f.FlightNumber = entity.OfficeVisit
			} else {
				f.FlightNumber = entity.Other
			}

			f.Status = entity.Await
			item.Flights = append(item.Flights, f)
		} else {
			for _, f := range flightRexExp.FindAllStringSubmatch(d.FlightNumber, -1) {
				if _, err := strconv.Atoi(f[1]); err != nil {
					log.Err(err).Msg(fmt.Sprintf("parse flight number (%s) error", f))
					continue
				}
				item.Flights = append(item.Flights, entity.Flight{
					FlightNumber: companyFlightPrefix + f[1],
					Airplane:     d.AircraftType,
					Status:       entity.Await,
				})
			}
		}

		o = append(o, item)
	}
	return o
}
