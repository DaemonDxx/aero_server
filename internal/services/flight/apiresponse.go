package flight

import (
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/api/flightaware"
)

func transformToEntity(r flightaware.Information) entity.FlightInfo {
	return entity.FlightInfo{
		FlightNumber:  r.FlightNumber,
		From:          r.From,
		To:            r.To,
		TimeDeparture: r.TimeDeparture,
		Duration:      r.Duration,
		AvgDuration:   r.AvgDuration,
	}
}
