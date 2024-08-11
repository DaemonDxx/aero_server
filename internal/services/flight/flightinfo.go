package flight

import (
	"context"
	"errors"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/api/flightaware"
	"github.com/daemondxx/lks_back/internal/dao"
	"github.com/daemondxx/lks_back/internal/services"
	"github.com/rs/zerolog"
	"strconv"
)

const servName = "flight_info_service"

type InfoDAO interface {
	GetByNumberFlight(ctx context.Context, number uint) (entity.FlightInfo, error)
	Save(ctx context.Context, info *entity.FlightInfo) error
}

type InfoAPI interface {
	GetFlightInfo(ctx context.Context, flight string) (flightaware.Information, error)
}

type InfoService struct {
	services.LoggedService
	d   InfoDAO
	api InfoAPI
}

func NewFlightInfoService(d InfoDAO, api InfoAPI, log *zerolog.Logger) *InfoService {
	l := log.With().Str("service", "flight_info_service").Logger()
	return &InfoService{
		d:             d,
		api:           api,
		LoggedService: services.NewLoggedService(&l),
	}
}

func (s *InfoService) GetFlightInfo(ctx context.Context, n uint) (entity.FlightInfo, error) {
	flight := "AFL" + strconv.Itoa(int(n))
	log := s.GetLogger("get_flight_info")
	log.Debug().Msg("find flight info in db..")
	info, err := s.d.GetByNumberFlight(ctx, n)

	if err == nil {
		log.Debug().Msg("flight info found from db")
		return info, nil
	}

	if !errors.Is(err, dao.ErrNotFoundInfo) {
		info.FlightNumber = flight
		return info, &services.ErrServ{
			Service: servName,
			Message: "find flight info in db error",
			Err:     err,
		}
	}
	log.Debug().Msg("flight info not found from db")

	log.Debug().Msg("get flight info from api...")
	r, err := s.api.GetFlightInfo(ctx, flight)
	if err != nil {
		info.FlightNumber = flight
		return info, &services.ErrServ{
			Service: servName,
			Message: "pull flight info from api error",
			Err:     err,
		}
	}

	info = transformToEntity(r)
	if err = s.d.Save(ctx, &info); err != nil {
		log.Err(err).Msg("save flight info in db error")
	}

	return info, nil
}
