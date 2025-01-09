package service_credential

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/api/lks"
	"github.com/daemondxx/lks_back/internal/logger"
	credential_mock "github.com/daemondxx/lks_back/mocks/server/services/credential"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type CreateCredentialSuite struct {
	suite.Suite
	dao     *credential_mock.MockCredentialDAO
	checker *credential_mock.MockAuthChecker
	accServ *credential_mock.MockAccountService
	serv    *Service
}

func (s *CreateCredentialSuite) SetupSuite() {
	s.dao = &credential_mock.MockCredentialDAO{}
	s.checker = &credential_mock.MockAuthChecker{}
	s.accServ = &credential_mock.MockAccountService{}
	s.serv = NewCredentialService(s.accServ, s.checker, s.dao, logger.NewLogger("DEV"))
}

func (s *CreateCredentialSuite) TestSuccessCreateWithoutFind() {
	var crID uint = 1
	cr := &entity.Credential{
		AccordLogin:    "acclogin",
		AccordPassword: "accPass",
		LKSLogin:       "lksLogin",
		LKSPassword:    "lksPass",
	}
	acc := &entity.Account{
		Model: gorm.Model{
			ID: 1,
		},
	}

	checkFn := s.checker.EXPECT().Check(mock.Anything, cr.AccordLogin, cr.AccordPassword, cr.LKSLogin, cr.LKSPassword).Return(nil)
	defer checkFn.Unset()
	findFn := s.dao.EXPECT().FindByLogin(mock.Anything, cr.AccordLogin, cr.LKSLogin).Return(nil, nil)
	defer findFn.Unset()
	createFn := s.dao.EXPECT().Create(mock.Anything, mock.Anything).Return(nil).Run(func(ctx context.Context, cred *entity.Credential) {
		cred.ID = crID
	})
	defer createFn.Unset()
	connFn := s.accServ.EXPECT().ConnectCredential(mock.Anything, acc, mock.Anything).Return(nil).Run(func(ctx context.Context, acc *entity.Account, cred *entity.Credential) {
		acc.CredentialID = &cred.ID
	})
	defer connFn.Unset()

	cr, err := s.serv.Create(context.Background(), acc, cr.AccordLogin, cr.AccordPassword, cr.LKSLogin, cr.LKSPassword)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), cr)
	assert.Equal(s.T(), crID, cr.ID)
	assert.Equal(s.T(), crID, *acc.CredentialID)
}

func (s *CreateCredentialSuite) TestSuccessCreateIfCreatedAndActual() {
	cr := &entity.Credential{
		Model: gorm.Model{
			ID: 1,
		},
		AccordLogin:    "acclogin",
		AccordPassword: "accPass",
		LKSLogin:       "lksLogin",
		LKSPassword:    "lksPass",
	}
	acc := &entity.Account{
		Model: gorm.Model{
			ID: 1,
		},
	}

	checkFn := s.checker.EXPECT().Check(mock.Anything, cr.AccordLogin, cr.AccordPassword, cr.LKSLogin, cr.LKSPassword).Return(nil)
	defer checkFn.Unset()
	findFn := s.dao.EXPECT().FindByLogin(mock.Anything, cr.AccordLogin, cr.LKSLogin).Return(cr, nil)
	defer findFn.Unset()
	connFn := s.accServ.EXPECT().ConnectCredential(mock.Anything, acc, mock.Anything).Return(nil).Run(func(ctx context.Context, acc *entity.Account, cred *entity.Credential) {
		acc.CredentialID = &cred.ID
	})
	defer connFn.Unset()

	cr, err := s.serv.Create(context.Background(), acc, cr.AccordLogin, cr.AccordPassword, cr.LKSLogin, cr.LKSPassword)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), cr)
	assert.Equal(s.T(), cr.ID, *acc.CredentialID)
}

func (s *CreateCredentialSuite) TestSuccessCreateIfCreatedAndNonActual() {
	cr := &entity.Credential{
		Model: gorm.Model{
			ID: 1,
		},
		AccordLogin:    "acclogin",
		AccordPassword: "accPass",
		LKSLogin:       "lksLogin",
		LKSPassword:    "lksPass",
		IsActual:       false,
	}
	actualAccPass := "newAccPass"
	actualLKSPass := "actualLKSPAss"

	acc := &entity.Account{
		Model: gorm.Model{
			ID: 1,
		},
	}

	checkFn := s.checker.EXPECT().Check(mock.Anything, cr.AccordLogin, actualAccPass, cr.LKSLogin, actualLKSPass).Return(nil)
	defer checkFn.Unset()
	findFn := s.dao.EXPECT().FindByLogin(mock.Anything, cr.AccordLogin, cr.LKSLogin).Return(cr, nil)
	defer findFn.Unset()
	saveFn := s.dao.EXPECT().Save(mock.Anything, cr).Return(nil)
	defer saveFn.Unset()
	connFn := s.accServ.EXPECT().ConnectCredential(mock.Anything, acc, mock.Anything).Return(nil).Run(func(ctx context.Context, acc *entity.Account, cred *entity.Credential) {
		acc.CredentialID = &cred.ID
	})
	defer connFn.Unset()

	cr, err := s.serv.Create(context.Background(), acc, cr.AccordLogin, actualAccPass, cr.LKSLogin, actualLKSPass)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), cr)
	assert.Equal(s.T(), cr.ID, *acc.CredentialID)
	assert.Equal(s.T(), cr.IsActual, true)
	assert.Equal(s.T(), actualAccPass, cr.AccordPassword)
	assert.Equal(s.T(), actualLKSPass, cr.LKSPassword)
}

func (s *CreateCredentialSuite) TestFailedCheck() {
	cr := &entity.Credential{
		AccordLogin:    "acclogin",
		AccordPassword: "accPass",
		LKSLogin:       "lksLogin",
		LKSPassword:    "lksPass",
	}

	acc := &entity.Account{
		Model: gorm.Model{
			ID: 1,
		},
	}

	checkFn := s.checker.EXPECT().Check(mock.Anything, cr.AccordLogin, cr.AccordPassword, cr.LKSLogin, cr.LKSPassword).Return(lks.ErrAccordAuth)
	defer checkFn.Unset()
	cr, err := s.serv.Create(context.Background(), acc, cr.AccordLogin, cr.AccordPassword, cr.LKSLogin, cr.LKSPassword)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, lks.ErrAccordAuth)
	assert.Nil(s.T(), acc.CredentialID)
}

func (s *CreateCredentialSuite) TestFailedCredentialIsConnected() {
	cr := &entity.Credential{
		AccordLogin:    "acclogin",
		AccordPassword: "accPass",
		LKSLogin:       "lksLogin",
		LKSPassword:    "lksPass",
	}

	var crID uint = 1
	acc := &entity.Account{
		Model: gorm.Model{
			ID: 1,
		},
		CredentialID: &crID,
	}

	cr, err := s.serv.Create(context.Background(), acc, cr.AccordLogin, cr.AccordPassword, cr.LKSLogin, cr.LKSPassword)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, ErrCredentialIsConnected)
	assert.Equal(s.T(), crID, *acc.CredentialID)
}

func TestService_CreateCredential(t *testing.T) {
	suite.Run(t, new(CreateCredentialSuite))
}
