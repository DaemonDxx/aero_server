package service_account

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/logger"
	"github.com/daemondxx/lks_back/internal/services"
	account_mock "github.com/daemondxx/lks_back/mocks/server/services/account"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CreateAccountSuite struct {
	suite.Suite
	dao       *account_mock.MockAccountDAO
	tokenServ *account_mock.MockTokenService
	service   *Service
}

func (s *CreateAccountSuite) BeforeTest(sName string, tName string) {
	l := logger.NewLogger("DEV")
	s.dao = account_mock.NewMockAccountDAO(s.T())
	s.tokenServ = account_mock.NewMockTokenService(s.T())
	s.service = &Service{
		LoggedService: services.NewLoggedService("test_account_service", l),
		dao:           s.dao,
		tokenServ:     s.tokenServ,
	}
}

func (s *CreateAccountSuite) TestSuccessCreate() {
	acc := &entity.Account{
		TelegramID: 1,
	}
	testToken := "test.token"

	findFn := s.dao.EXPECT().Find(mock.Anything, acc.TelegramID).Return(nil, nil)
	defer findFn.Unset()

	createFn := s.dao.EXPECT().Create(mock.Anything, acc).Return(nil)
	defer createFn.Unset()

	createTFn := s.tokenServ.EXPECT().Create(acc).Return(testToken)
	defer createTFn.Unset()

	t, err := s.service.CreateAccount(context.Background(), acc.TelegramID)

	require.NoError(s.T(), err)
	require.Equal(s.T(), testToken, t)
}

func (s *CreateAccountSuite) TestErrAccIsExist() {
	acc := &entity.Account{
		TelegramID: 1,
	}

	findFn := s.dao.EXPECT().Find(mock.Anything, acc.TelegramID).Return([]entity.Account{{TelegramID: 1}}, nil)
	defer findFn.Unset()

	t, err := s.service.CreateAccount(context.Background(), acc.TelegramID)

	require.ErrorIs(s.T(), err, ErrAccountExists)
	require.Equal(s.T(), "", t)
}

func (s *CreateAccountSuite) TestErrDAOFind() {
	acc := &entity.Account{
		TelegramID: 1,
	}

	findFn := s.dao.EXPECT().Find(mock.Anything, acc.TelegramID).Return(nil, errors.New("something wrong"))
	defer findFn.Unset()

	t, err := s.service.CreateAccount(context.Background(), acc.TelegramID)

	var expectedErr *services.ErrServ
	require.ErrorAs(s.T(), err, &expectedErr)
	require.Equal(s.T(), "", t)
}

func (s *CreateAccountSuite) TestErrDAOCreate() {
	acc := &entity.Account{
		TelegramID: 1,
	}

	findFn := s.dao.EXPECT().Find(mock.Anything, acc.TelegramID).Return(nil, nil)
	defer findFn.Unset()

	createFn := s.dao.EXPECT().Create(mock.Anything, acc).Return(errors.New("something wrong"))
	defer createFn.Unset()

	t, err := s.service.CreateAccount(context.Background(), acc.TelegramID)

	var expectedErr *services.ErrServ
	require.ErrorAs(s.T(), err, &expectedErr)
	require.Equal(s.T(), "", t)
}

func TestCreateAccountSuite(t *testing.T) {
	suite.Run(t, new(CreateAccountSuite))
}
