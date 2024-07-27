package lks

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	flightNumberAttr = "FlightNumber"
	airplaneAttr     = "AircraftType"
	departureAttr    = "StartDate"
	arrivalAttr      = "EndDate"
	noteAttr         = "Note"
	routeAttr        = "Route"
	targetAttr       = "Target"
	plateAttr        = "Place"
	confirmedAttr    = "ConfirmTypeStr()"
)

var confirmedTimeRegExp, _ = regexp.Compile("(\\d\\d.\\d\\d.\\d\\d\\d\\d \\d\\d:\\d\\d)")

type CurrentOrderRow struct {
	Flights   []uint
	Airplane  string
	Departure time.Time
	Arrival   time.Time
	Target    string
	Place     string
	Note      string
	Route     string
	Confirmed *time.Time
}

func newCOrderRow(m map[string]string) (CurrentOrderRow, error) {
	row := CurrentOrderRow{}

	t := m[flightNumberAttr]
	if t == "" {
		//todo подумать над отпусками и нарядами
	} else {
		t := strings.Split(t, " ")
		for _, fStr := range t {
			fNumb, err := strconv.Atoi(fStr)
			if err != nil {
				return row, fmt.Errorf("parse flight number (%s) error: %e", fStr, err)
			}
			row.Flights = append(row.Flights, uint(fNumb))
		}
	}

	dTime, err := parseLKSFormatData(m[departureAttr])
	if err != nil {
		return row, fmt.Errorf("parse departure time (%s) error: %e", m[departureAttr], err)
	}
	aTime, err := parseLKSFormatData(m[arrivalAttr])
	if err != nil {
		return row, fmt.Errorf("parse arrival time (%s) error: %e", m[arrivalAttr], err)
	}
	row.Departure = dTime
	row.Arrival = aTime

	row.Airplane = m[airplaneAttr]
	row.Note = m[noteAttr]
	row.Place = m[plateAttr]
	row.Target = m[targetAttr]
	row.Route = m[routeAttr]

	if m[confirmedAttr] != "" {
		r := confirmedTimeRegExp.FindStringSubmatch(m[confirmedAttr])
		if len(r) == 2 {
			conTime, err := parseLKSConfirmedFormatData(r[1])
			if err != nil {
				return row, fmt.Errorf("parse confirmed time (%s) error: %e", m[confirmedAttr], err)
			}
			row.Confirmed = &conTime
		}
	}
	return row, nil
}
