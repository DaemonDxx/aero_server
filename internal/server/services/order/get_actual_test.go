package service_order

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/logger"
	"github.com/daemondxx/lks_back/internal/services"
	service_order_mock "github.com/daemondxx/lks_back/mocks/server/services/order"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type GetActualOrderSuite struct {
	suite.Suite
	oDao  *service_order_mock.MockOrderDAO
	crDAO *service_order_mock.MockCredentialDAO
	serv  *Service
	acc   *entity.Account
}

func TestService_GetActualOrders(t *testing.T) {
	suite.Run(t, new(GetActualOrderSuite))
}

func (s *GetActualOrderSuite) SetupSuite() {
	log := logger.NewLogger("DEV")
	s.oDao = &service_order_mock.MockOrderDAO{}
	s.crDAO = &service_order_mock.MockCredentialDAO{}

	var credID uint = 1
	s.crDAO.EXPECT().GetByID(mock.Anything, credID).Return(&entity.Credential{Model: gorm.Model{ID: credID}}, nil)
	s.serv = &Service{
		LoggedService: services.NewLoggedService("order_test", log),
		orderDAO:      s.oDao,
		crDAO:         s.crDAO,
	}

	s.acc = &entity.Account{
		Model: gorm.Model{
			ID: 1,
		},
		TelegramID:   0,
		CredentialID: &credID,
	}
}

func (s *GetActualOrderSuite) TestReturnEmptyOrders() {
	findFn := s.oDao.EXPECT().FindLastOrders(mock.Anything, mock.Anything, 2).Return([]entity.Order{}, nil)
	defer findFn.Unset()

	o, err := s.serv.GetActualOrders(context.Background(), *s.acc)

	require.NoError(s.T(), err)
	require.Equal(s.T(), 0, len(o))
}

func (s *GetActualOrderSuite) TestReturnTwoActualOrders() {
	findFn := s.oDao.EXPECT().FindLastOrders(mock.Anything, mock.Anything, 2).Return(makeActualOrders(2), nil)
	defer findFn.Unset()

	o, err := s.serv.GetActualOrders(context.Background(), *s.acc)

	require.NoError(s.T(), err)
	require.Equal(s.T(), 2, len(o))
}

func (s *GetActualOrderSuite) TestReturnOneActualOrder() {
	var o []entity.Order

	o = append(o, makeActualOrders(1)...)
	o = append(o, makeNonActualOrders(1)...)

	findFn := s.oDao.EXPECT().FindLastOrders(mock.Anything, mock.Anything, 2).Return(o, nil)
	defer findFn.Unset()

	o, err := s.serv.GetActualOrders(context.Background(), *s.acc)

	require.NoError(s.T(), err)
	require.Equal(s.T(), 1, len(o))
}

func (s *GetActualOrderSuite) TestReturnOneActualOrderIfHasEmptyOrder() {
	var o []entity.Order

	o = append(o, makeActualOrders(1)...)
	o = append(o, entity.Order{
		Model:        gorm.Model{},
		CredentialID: 0,
		Items:        make([]entity.OrderItem, 0),
		Status:       0,
	})

	findFn := s.oDao.EXPECT().FindLastOrders(mock.Anything, mock.Anything, 2).Return(o, nil)
	defer findFn.Unset()

	o, err := s.serv.GetActualOrders(context.Background(), *s.acc)

	require.NoError(s.T(), err)
	require.Equal(s.T(), 1, len(o))
}

func makeActualOrders(count int) []entity.Order {
	o := make([]entity.Order, 0, count)
	for i := 0; i < count; i++ {
		o = append(o, entity.Order{
			Model: gorm.Model{
				ID: uint(count - i),
			},
			CredentialID: 1,
			Items: []entity.OrderItem{
				{
					Departure:   time.Now().Add(time.Duration(count-i) * time.Hour),
					Description: "Test",
					Route:       "Test",
					ConfirmDate: nil,
					OrderID:     1,
				},
			},
			Status: 0,
		})
	}
	return o
}

func makeNonActualOrders(count int) []entity.Order {
	o := make([]entity.Order, 0, count)
	for i := 0; i < count; i++ {
		o = append(o, entity.Order{
			Model: gorm.Model{
				ID: uint(count - i),
			},
			CredentialID: 1,
			Items: []entity.OrderItem{
				{
					Departure:   time.Now().Add(-time.Duration(count-i) * time.Hour),
					Description: "Test",
					Route:       "Test",
					ConfirmDate: nil,
					OrderID:     1,
				},
			},
			Status: 0,
		})
	}
	return o
}
