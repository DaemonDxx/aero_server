package order

import (
	"context"
	"fmt"
	entity "github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/api/lks"
	"github.com/daemondxx/lks_back/internal/dao"
	"github.com/daemondxx/lks_back/internal/services"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"regexp"
	"strconv"
	"strings"
)

const servName = "order_service"
const companyFlightPrefix = "AFL"

var flightRexExp, _ = regexp.Compile("(\\d+)")

var errHaveNotNewItems = errors.New("slice have not new items")

type Service struct {
	services.LoggedService
	oDAO  DAO
	lks   LksAPI
	uServ UserService
}

func NewOrderService(dao DAO, l LksAPI, u UserService, log *zerolog.Logger) *Service {
	logger := log.With().Str("service", "order_service").Logger()
	return &Service{
		LoggedService: services.NewLoggedService(&logger),
		oDAO:          dao,
		lks:           l,
		uServ:         u,
	}
}

func (s *Service) GetActualOrder(ctx context.Context, userID uint) (entity.Order, error) {
	l := s.GetLogger("get_actual_order").With().Uint("user_id", userID).Logger()
	var order entity.Order

	u, err := s.uServ.GetUserByID(ctx, userID)
	if err != nil {
		l.Err(err).Msg("get user info error")
		return order, &services.ErrServ{
			Service: servName,
			Message: "get auth info user error",
			Err:     err,
		}
	}

	aDuty, err := s.lks.GetActualDuty(ctx, lks.AuthPayload{
		AccordLogin:    u.AccordLogin,
		AccordPassword: u.AccordPassword,
		LksLogin:       u.LKSLogin,
		LksPassword:    u.LKSPassword,
	})
	if err != nil {
		l.Err(err).Msg("get actual order from api error")
		return entity.Order{}, &services.ErrServ{
			Service: servName,
			Message: "get actual order from api error",
			Err:     err,
		}
	}

	if len(aDuty) == 0 {
		return newEmptyOrder(userID), nil
	}

	items := s.transformToOrderItems(aDuty)

	last, err := s.oDAO.GetLastOrder(ctx, userID)
	if err != nil && !errors.Is(err, dao.ErrOrderNotFound) {
		return order, &services.ErrServ{
			Service: servName,
			Message: "get last order error",
			Err:     err,
		}
	}

	newItems, err := s.extractNewItems(last, items)

	//create new order
	if err == nil {
		order.UserID = userID
		order.Items = newItems

		if newItems[0].ConfirmDate != nil {
			order.Status = entity.Confirm
		} else {
			order.Status = entity.AwaitConfirmation
		}

		if err := s.oDAO.Create(ctx, &order); err != nil {
			l.Err(err).Msg("save new order error")
		}

		return order, nil
	}

	//unknown error
	if !errors.Is(err, errHaveNotNewItems) {
		log.Err(err).Msg("extract new items by duty error")
		return order, &services.ErrServ{
			Service: servName,
			Message: "extract new items by duty error",
			Err:     err,
		}
	}

	//update status by last order
	if len(last.Items) == 0 || last.Status == entity.Limited {
		return last, nil
	}

	if aDuty[len(aDuty)-1].BlockDate != nil {
		last.Status = entity.Limited
		if err := s.oDAO.Save(ctx, &last); err != nil {
			l.Err(err).Msg("change status on limited error")
		}
		return last, nil
	} else if last.Status != entity.Confirm && aDuty[len(aDuty)-1].ConfirmDate != nil {
		last.Status = entity.Confirm
		if err := s.oDAO.Save(ctx, &last); err != nil {
			l.Err(err).Msg("change status on confirm error")
		}
		return last, nil
	}

	return last, nil
}

func (s *Service) transformToOrderItems(aDuty []lks.CurrentDuty) []entity.OrderItem {
	items := make([]entity.OrderItem, 0, len(aDuty))
	for _, d := range aDuty {
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
				f.FlightNumber = DayReserve
			} else if strings.Contains(d.Note, "ночной") {
				f.FlightNumber = NightReserve
			} else if strings.Contains(d.Note, "Отпуск") {
				f.FlightNumber = Holyday
			} else if strings.Contains(d.Note, "Явка") {
				f.FlightNumber = OfficeVisit
			} else {
				f.FlightNumber = Other
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

		items = append(items, item)
	}
	return items
}

// items length dont be equal 0
// Can return ErrNotFound
func (s *Service) extractNewItems(last entity.Order, items []entity.OrderItem) ([]entity.OrderItem, error) {
	if last.ID == 0 {
		res := make([]entity.OrderItem, len(items))
		copy(res, items)
		return res, nil
	}

	lItem := last.Items[len(last.Items)-1]
	if items[len(items)-1].Departure.Unix() <= lItem.Departure.Unix() {
		return nil, errHaveNotNewItems
	}

	for i := len(items) - 1; i >= 0; i-- {
		if items[i].Departure.Unix() <= lItem.Departure.Unix() {
			items = items[i+1:]
			res := make([]entity.OrderItem, len(items))
			copy(res, items)
			return res, nil
		}
	}

	res := make([]entity.OrderItem, len(items))
	copy(res, items)
	return res, nil
}
