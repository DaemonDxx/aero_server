package puller

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/logger"
	"github.com/daemondxx/lks_back/internal/services"
	puller_mock "github.com/daemondxx/lks_back/mocks/server/services/puller"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type collectorPullSuite struct {
	suite.Suite
	dao  *puller_mock.MockOrderDAO
	api  *puller_mock.MockLKSApi
	cr   *entity.Credential
	serv *Service
}

func TestCollectorPullSuite(t *testing.T) {
	suite.Run(t, new(collectorPullSuite))
}

func (s *collectorPullSuite) SetupTest() {
	s.dao = &puller_mock.MockOrderDAO{}
	s.api = &puller_mock.MockLKSApi{}
	s.cr = &entity.Credential{
		Model: gorm.Model{
			ID: 1,
		},
		AccordLogin:    "Test",
		LKSLogin:       "Test",
		AccordPassword: "Test",
		LKSPassword:    "Test",
		IsActual:       true,
	}
	s.serv = &Service{
		LoggedService: services.NewLoggedService("collector_test", logger.NewLogger("DEV")),
		dao:           s.dao,
		api:           s.api,
	}
}

func (s *collectorPullSuite) TestSuccessPullFirstOrder() {
	pi := []entity.OrderItem{
		{
			Departure: getTimeWithOffset(5),
		},
	}
	findFn := s.dao.EXPECT().FindLastOrders(mock.Anything, s.cr, 1).Return(make([]entity.Order, 0), nil)
	defer findFn.Unset()
	apiFn := s.api.EXPECT().GetActualDuty(mock.Anything, s.cr).Return(pi, nil)
	defer apiFn.Unset()
	saveFn := s.dao.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)
	defer saveFn.Unset()

	o, err := s.serv.Pull(context.Background(), s.cr)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), o)

	assert.Equal(s.T(), len(pi), len(o.Items))
	assert.Equal(s.T(), s.cr.ID, o.CredentialID)
}

func (s *collectorPullSuite) TestPullOldOrder() {
	pi := []entity.OrderItem{
		{
			Departure: getTimeWithOffset(5),
		},
	}
	o := entity.Order{
		Model:        gorm.Model{},
		CredentialID: s.cr.ID,
		Items: []entity.OrderItem{
			{
				Departure: getTimeWithOffset(5),
			},
		},
		Status: 0,
	}
	findFn := s.dao.EXPECT().FindLastOrders(mock.Anything, s.cr, 1).Return([]entity.Order{o}, nil)
	defer findFn.Unset()
	apiFn := s.api.EXPECT().GetActualDuty(mock.Anything, s.cr).Return(pi, nil)
	defer apiFn.Unset()

	or, err := s.serv.Pull(context.Background(), s.cr)
	require.ErrorIs(s.T(), err, ErrOrderIsExist)
	require.Nil(s.T(), or)
}

func (s *collectorPullSuite) TestPullNewOrderWithEmptyLast() {
	pi := []entity.OrderItem{
		{
			Departure: getTimeWithOffset(5),
		},
	}
	o := entity.Order{
		Model:        gorm.Model{},
		CredentialID: s.cr.ID,
		Items:        []entity.OrderItem{},
		Status:       0,
	}
	findFn := s.dao.EXPECT().FindLastOrders(mock.Anything, s.cr, 1).Return([]entity.Order{o}, nil)
	defer findFn.Unset()
	apiFn := s.api.EXPECT().GetActualDuty(mock.Anything, s.cr).Return(pi, nil)
	defer apiFn.Unset()
	saveFn := s.dao.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)
	defer saveFn.Unset()

	or, err := s.serv.Pull(context.Background(), s.cr)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), or)

	assert.Equal(s.T(), len(pi), len(or.Items))
	assert.Equal(s.T(), s.cr.ID, or.CredentialID)
}

func (s *collectorPullSuite) TestPullNewOrderWithLast() {
	pi := []entity.OrderItem{
		{
			Departure: getTimeWithOffset(5),
		},
		{
			Departure: getTimeWithOffset(7),
		},
	}
	o := entity.Order{
		Model:        gorm.Model{},
		CredentialID: s.cr.ID,
		Items: []entity.OrderItem{
			{
				Departure: getTimeWithOffset(5),
			},
		},
		Status: 0,
	}
	findFn := s.dao.EXPECT().FindLastOrders(mock.Anything, s.cr, 1).Return([]entity.Order{o}, nil)
	defer findFn.Unset()
	apiFn := s.api.EXPECT().GetActualDuty(mock.Anything, s.cr).Return(pi, nil)
	defer apiFn.Unset()
	saveFn := s.dao.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)
	defer saveFn.Unset()

	or, err := s.serv.Pull(context.Background(), s.cr)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), or)

	assert.Equal(s.T(), 1, len(or.Items))
	assert.Equal(s.T(), s.cr.ID, or.CredentialID)
}

func (s *collectorPullSuite) TestPullEmptyOrder() {
	var pi []entity.OrderItem
	o := entity.Order{
		Model:        gorm.Model{},
		CredentialID: s.cr.ID,
		Items: []entity.OrderItem{
			{
				Departure: getTimeWithOffset(5),
			},
		},
		Status: 0,
	}
	findFn := s.dao.EXPECT().FindLastOrders(mock.Anything, s.cr, 1).Return([]entity.Order{o}, nil)
	defer findFn.Unset()
	apiFn := s.api.EXPECT().GetActualDuty(mock.Anything, s.cr).Return(pi, nil)
	defer apiFn.Unset()
	saveFn := s.dao.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)
	defer saveFn.Unset()

	or, err := s.serv.Pull(context.Background(), s.cr)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), or)

	assert.Equal(s.T(), 0, len(or.Items))
	assert.Equal(s.T(), s.cr.ID, or.CredentialID)
}
