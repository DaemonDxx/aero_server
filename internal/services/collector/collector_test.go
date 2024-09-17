package collector

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/logger"
	"github.com/daemondxx/lks_back/internal/services"
	m "github.com/daemondxx/lks_back/tests/mocks/collector"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type TestCollectActualOrderSuite struct {
	suite.Suite
	oServ *m.MockOrderService
	uDao  *m.MockUserDAO
	nServ *m.MockNotificationService
	log   *zerolog.Logger

	fnFind          *m.MockUserDAO_Find_Call
	fnGet           *m.MockOrderService_GetActualOrder_Call
	fnSuccessNotify *m.MockNotificationService_ActualOrderNotify_Call
	fnErrNotify     *m.MockNotificationService_ErrorNotify_Call
}

func (s *TestCollectActualOrderSuite) SetupSuite() {
	s.log = logger.NewLogger(logger.DEV)
	s.oServ = &m.MockOrderService{}
	s.uDao = &m.MockUserDAO{}
	s.nServ = &m.MockNotificationService{}
}

func (s *TestCollectActualOrderSuite) AfterTest(sName string, tName string) {
	s.fnFind.Unset()
	s.fnGet.Unset()
	s.fnSuccessNotify.Unset()
	s.fnErrNotify.Unset()
}

func (s *TestCollectActualOrderSuite) TestSuccessCollect() {
	users := []entity.User{
		{
			ID: 0,
		},
	}

	order := entity.Order{
		UserID: users[0].ID,
	}

	c := NewCollectorService(s.uDao, s.oServ, s.nServ, Config{
		MaxAttempts:     1,
		MinTimeoutRetry: 100 * time.Millisecond,
	}, s.log)

	s.fnFind = s.uDao.EXPECT().Find(mock.Anything, mock.Anything).Return(users, nil)
	s.fnGet = s.oServ.EXPECT().GetActualOrder(mock.Anything, users[0].ID).Return(order, nil)
	s.fnSuccessNotify = s.nServ.EXPECT().ActualOrderNotify(users[0].ID, order).Return()
	s.fnErrNotify = s.nServ.EXPECT().ErrorNotify(mock.Anything, mock.Anything).Return()

	err := c.CollectActualOrder(context.Background())

	s.fnFind.Once()
	s.fnGet.Once()
	s.fnSuccessNotify.Once()
	s.fnErrNotify.Times(0)

	assert.NoError(s.T(), err)
}

func (s *TestCollectActualOrderSuite) TestLimitAttemptErrorCollect() {
	users := []entity.User{
		{
			ID: 0,
		},
		{
			ID: 2,
		},
	}

	maxAttempts := 3

	c := NewCollectorService(s.uDao, s.oServ, s.nServ, Config{
		MaxAttempts:     uint(maxAttempts),
		MinTimeoutRetry: 100 * time.Microsecond,
	}, s.log)

	s.fnFind = s.uDao.EXPECT().Find(mock.Anything, mock.Anything).Return(users, nil)
	s.fnGet = s.oServ.EXPECT().GetActualOrder(mock.Anything, mock.Anything).Return(entity.Order{}, context.DeadlineExceeded)
	s.fnSuccessNotify = s.nServ.EXPECT().ActualOrderNotify(mock.Anything, mock.Anything).Return()
	s.fnErrNotify = s.nServ.EXPECT().ErrorNotify(mock.Anything, mock.Anything).Return()

	err := c.CollectActualOrder(context.Background())

	s.fnFind.Once()
	s.fnGet.Times(maxAttempts * len(users))
	s.fnSuccessNotify.Times(0)
	s.fnErrNotify.Times(len(users))

	var target *ErrLimitAttempt
	require.ErrorAs(s.T(), err, &target)
	e := err.(*services.ErrServ).Err.(*ErrLimitAttempt)
	assert.Equal(s.T(), len(users), len(e.Users))
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

	var users []entity.User

	for k, _ := range cases {
		users = append(users, entity.User{ID: k})
	}

	maxAttempts := 5

	c := NewCollectorService(s.uDao, s.oServ, s.nServ, Config{
		MaxAttempts:     uint(maxAttempts),
		MinTimeoutRetry: 10 * time.Microsecond,
	}, s.log)

	s.fnFind = s.uDao.EXPECT().Find(mock.Anything, mock.Anything).Return(users, nil)
	s.fnGet = s.oServ.EXPECT().GetActualOrder(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, u uint) (entity.Order, error) {
		r, _ := cases[u]
		if r.attempt+1 >= r.attemptBySuccessful {
			return entity.Order{}, nil
		} else {
			r.attempt++
			return entity.Order{}, context.DeadlineExceeded
		}
	})
	s.fnSuccessNotify = s.nServ.EXPECT().ActualOrderNotify(mock.Anything, mock.Anything).Return()
	s.fnErrNotify = s.nServ.EXPECT().ErrorNotify(mock.Anything, mock.Anything).Return()

	err := c.CollectActualOrder(context.Background())

	s.fnFind.Once()
	s.fnGet.Times(6)
	s.fnSuccessNotify.Times(len(users))
	s.fnErrNotify.Times(0)

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

	var users []entity.User

	for k, _ := range cases {
		users = append(users, entity.User{ID: k})
	}

	maxAttempts := 2

	c := NewCollectorService(s.uDao, s.oServ, s.nServ, Config{
		MaxAttempts:     uint(maxAttempts),
		MinTimeoutRetry: 10 * time.Microsecond,
	}, s.log)

	s.fnFind = s.uDao.EXPECT().Find(mock.Anything, mock.Anything).Return(users, nil)
	s.fnGet = s.oServ.EXPECT().GetActualOrder(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, u uint) (entity.Order, error) {
		r, _ := cases[u]
		if r.attempt+1 >= r.attemptBySuccessful {
			return entity.Order{}, nil
		} else {
			r.attempt++
			return entity.Order{}, context.DeadlineExceeded
		}
	})
	s.fnSuccessNotify = s.nServ.EXPECT().ActualOrderNotify(mock.Anything, mock.Anything).Return()
	s.fnErrNotify = s.nServ.EXPECT().ErrorNotify(mock.Anything, mock.Anything).Return()

	err := c.CollectActualOrder(context.Background())

	s.fnFind.Once()
	s.fnGet.Times(5)
	s.fnSuccessNotify.Times(len(users) - 1)
	s.fnErrNotify.Times(1)

	var target *ErrLimitAttempt
	require.ErrorAs(s.T(), err, &target)
	e := err.(*services.ErrServ).Err.(*ErrLimitAttempt)
	assert.Equal(s.T(), 1, len(e.Users))
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

	var users []entity.User

	for k, _ := range cases {
		users = append(users, entity.User{ID: k})
	}

	maxAttempts := 3
	timeout := 500 * time.Millisecond

	c := NewCollectorService(s.uDao, s.oServ, s.nServ, Config{
		MaxAttempts:     uint(maxAttempts),
		MinTimeoutRetry: timeout,
	}, s.log)

	s.fnFind = s.uDao.EXPECT().Find(mock.Anything, mock.Anything).Return(users, nil)
	s.fnGet = s.oServ.EXPECT().GetActualOrder(mock.Anything, mock.Anything).RunAndReturn(func(ctx context.Context, u uint) (entity.Order, error) {
		r, _ := cases[u]
		if r.attempt+1 >= r.attemptBySuccessful {
			return entity.Order{}, nil
		} else {
			r.attempt++
			return entity.Order{}, context.DeadlineExceeded
		}
	})
	s.fnSuccessNotify = s.nServ.EXPECT().ActualOrderNotify(mock.Anything, mock.Anything).Return()
	s.fnErrNotify = s.nServ.EXPECT().ErrorNotify(mock.Anything, mock.Anything).Return()

	timeStart := time.Now()
	err := c.CollectActualOrder(context.Background())

	d := time.Now().Sub(timeStart)

	s.fnFind.Once()
	s.fnGet.Times(3)
	s.fnSuccessNotify.Times(1)
	s.fnErrNotify.Times(0)

	assert.Greater(s.T(), int64(d), int64(timeout)*int64(maxAttempts-1))
	assert.NoError(s.T(), err)
}

func TestCollectorActualOrderSuite(t *testing.T) {
	suite.Run(t, new(TestCollectActualOrderSuite))
}
