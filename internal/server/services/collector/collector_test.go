package collector

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/logger"
	m "github.com/daemondxx/lks_back/mocks/server/services/collector"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type TestCollectActualOrderSuite struct {
	suite.Suite
	oServ *m.MockOrderService
	uDao  *m.MockCredentialDAO
	nServ *m.MockNotificationService
	log   *zerolog.Logger

	fnFind   *m.MockCredentialDAO_GetActualCredential_Call
	fnPull   *m.MockOrderService_PullNewOrder_Call
	fnNotify *m.MockNotificationService_Send_Call
}

func (s *TestCollectActualOrderSuite) wait() {
	t := time.NewTimer(5 * time.Millisecond)
	<-t.C
}

func (s *TestCollectActualOrderSuite) SetupSuite() {
	s.log = logger.NewLogger(logger.DEV)
	s.oServ = &m.MockOrderService{}
	s.uDao = &m.MockCredentialDAO{}
	s.nServ = &m.MockNotificationService{}
}

func (s *TestCollectActualOrderSuite) AfterTest(sName string, tName string) {
	s.fnFind.Unset()
	s.fnPull.Unset()
	s.fnNotify.Unset()
}

func (s *TestCollectActualOrderSuite) TestSuccessCollect() {
	crs := []entity.Credential{
		{
			Model: gorm.Model{ID: 0},
		},
	}

	order := entity.Order{
		CredentialID: crs[0].ID,
	}

	c := NewCollectorService(s.uDao, s.oServ, s.nServ, Config{
		MaxAttempts:     1,
		MinTimeoutRetry: 100 * time.Millisecond,
	}, s.log)

	s.fnFind = s.uDao.EXPECT().GetActualCredential(mock.Anything).Return(crs, nil)
	s.fnPull = s.oServ.EXPECT().PullNewOrder(mock.Anything, &crs[0]).Return(&order, nil)
	s.fnNotify = s.nServ.EXPECT().Send(mock.Anything).Return(nil)

	err := c.CollectActualOrder(context.Background())
	s.wait()

	s.fnFind.Once()
	s.fnPull.Once()
	s.fnNotify.Once()

	assert.NoError(s.T(), err)
}

func (s *TestCollectActualOrderSuite) TestLimitAttemptErrorCollect() {
	cr := []entity.Credential{
		{
			Model: gorm.Model{ID: 0},
		},
		{
			Model: gorm.Model{ID: 2},
		},
	}

	maxAttempts := 3

	c := NewCollectorService(s.uDao, s.oServ, s.nServ, Config{
		MaxAttempts:     maxAttempts,
		MinTimeoutRetry: 100 * time.Microsecond,
	}, s.log)

	s.fnFind = s.uDao.EXPECT().GetActualCredential(mock.Anything).Return(cr, nil)
	s.fnPull = s.oServ.EXPECT().PullNewOrder(mock.Anything, mock.Anything).Return(&entity.Order{}, context.DeadlineExceeded)
	s.fnNotify = s.nServ.EXPECT().Send(mock.Anything).Return(nil)

	err := c.CollectActualOrder(context.Background())
	s.wait()

	s.fnFind.Once()
	s.fnPull.Times(maxAttempts * len(cr))
	s.fnNotify.Times(maxAttempts * len(cr))

	var target *ErrLimitAttempt
	assert.ErrorAs(s.T(), err, &target)
	assert.Equal(s.T(), len(cr), len(target.Credentials))
}

func (s *TestCollectActualOrderSuite) TestRetryWithAllSuccessfulResultCollect() {

	type orderCase struct {
		attemptBySuccessful int
		attempt             int
	}

	cases := make(map[uint]*orderCase)
	cases[0] = &orderCase{
		attemptBySuccessful: 1,
	}
	cases[1] = &orderCase{
		attemptBySuccessful: 2,
	}
	cases[2] = &orderCase{
		attemptBySuccessful: 3,
	}

	var crs []entity.Credential

	for k, _ := range cases {
		crs = append(crs, entity.Credential{Model: gorm.Model{ID: k}})
	}

	maxAttempts := 5

	c := NewCollectorService(s.uDao, s.oServ, s.nServ, Config{
		MaxAttempts:     maxAttempts,
		MinTimeoutRetry: 10 * time.Microsecond,
	}, s.log)

	s.fnFind = s.uDao.EXPECT().GetActualCredential(mock.Anything).Return(crs, nil)
	s.fnPull = s.oServ.EXPECT().PullNewOrder(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, cr *entity.Credential) (*entity.Order, error) {
		r, _ := cases[cr.ID]
		if r.attempt+1 >= r.attemptBySuccessful {
			return &entity.Order{}, nil
		} else {
			r.attempt++
			return &entity.Order{}, context.DeadlineExceeded
		}
	})
	s.fnNotify = s.nServ.EXPECT().Send(mock.Anything).Return(nil)

	err := c.CollectActualOrder(context.Background())
	s.wait()

	s.fnFind.Once()
	s.fnPull.Times(6)
	s.fnNotify.Times(len(crs))

	assert.NoError(s.T(), err)
}

func (s *TestCollectActualOrderSuite) TestRetryWithSomeFailureResultCollect() {

	type orderCase struct {
		attemptBySuccessful int
		attempt             int
	}

	cases := make(map[uint]*orderCase)
	cases[0] = &orderCase{
		attemptBySuccessful: 1,
	}
	cases[1] = &orderCase{
		attemptBySuccessful: 2,
	}
	cases[2] = &orderCase{
		attemptBySuccessful: 3,
	}

	var crs []entity.Credential

	for k, _ := range cases {
		crs = append(crs, entity.Credential{Model: gorm.Model{ID: k}})
	}

	maxAttempts := 2

	c := NewCollectorService(s.uDao, s.oServ, s.nServ, Config{
		MaxAttempts:     maxAttempts,
		MinTimeoutRetry: 10 * time.Microsecond,
	}, s.log)

	s.fnFind = s.uDao.EXPECT().GetActualCredential(mock.Anything).Return(crs, nil)
	s.fnPull = s.oServ.EXPECT().PullNewOrder(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, cr *entity.Credential) (*entity.Order, error) {
		r, _ := cases[cr.ID]
		if r.attempt+1 >= r.attemptBySuccessful {
			return &entity.Order{}, nil
		} else {
			r.attempt++
			return &entity.Order{}, context.DeadlineExceeded
		}
	})
	s.fnNotify = s.nServ.EXPECT().Send(mock.Anything).Return(nil)

	err := c.CollectActualOrder(context.Background())
	s.wait()

	s.fnFind.Once()
	s.fnPull.Times(5)
	s.fnNotify.Times(len(crs) - 1)

	var target *ErrLimitAttempt
	assert.ErrorAs(s.T(), err, &target)
	assert.Equal(s.T(), 1, len(target.Credentials))
}

func (s *TestCollectActualOrderSuite) TestTimeoutCollectCollect() {

	type orderCase struct {
		attemptBySuccessful int
		attempt             int
	}

	cases := make(map[uint]*orderCase)
	cases[0] = &orderCase{
		attemptBySuccessful: 3,
	}

	var crs []entity.Credential

	for k, _ := range cases {
		crs = append(crs, entity.Credential{Model: gorm.Model{ID: k}})
	}

	maxAttempts := 3
	timeout := 500 * time.Millisecond

	c := NewCollectorService(s.uDao, s.oServ, s.nServ, Config{
		MaxAttempts:     maxAttempts,
		MinTimeoutRetry: timeout,
	}, s.log)

	s.fnFind = s.uDao.EXPECT().GetActualCredential(mock.Anything).Return(crs, nil)
	s.fnPull = s.oServ.EXPECT().PullNewOrder(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, cr *entity.Credential) (*entity.Order, error) {
		r, _ := cases[cr.ID]
		if r.attempt+1 >= r.attemptBySuccessful {
			return &entity.Order{}, nil
		} else {
			r.attempt++
			return &entity.Order{}, context.DeadlineExceeded
		}
	})
	s.fnNotify = s.nServ.EXPECT().Send(mock.Anything).Return(nil)

	timeStart := time.Now()
	err := c.CollectActualOrder(context.Background())
	s.wait()

	d := time.Now().Sub(timeStart)

	s.fnFind.Once()
	s.fnPull.Times(3)
	s.fnNotify.Times(1)

	assert.Greater(s.T(), int64(d), int64(timeout)*int64(maxAttempts-1))
	assert.NoError(s.T(), err)
}

func TestCollectorActualOrderSuite(t *testing.T) {
	suite.Run(t, new(TestCollectActualOrderSuite))
}
