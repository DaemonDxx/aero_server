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

type UpdateCredentialsSuite struct {
	suite.Suite
	dao     *credential_mock.MockCredentialDAO
	checker *credential_mock.MockAuthChecker
	serv    *Service
}

func (s *UpdateCredentialsSuite) SetupSuite() {
	s.dao = &credential_mock.MockCredentialDAO{}
	s.checker = &credential_mock.MockAuthChecker{}
	accServ := &credential_mock.MockAccountService{}
	s.serv = NewCredentialService(accServ, s.checker, s.dao, logger.NewLogger("DEV"))
}

func (s *UpdateCredentialsSuite) TestSuccessUpdate() {
	cr := &entity.Credential{
		Model: gorm.Model{
			ID: 1,
		},
		AccordLogin:    "acclogin",
		LKSLogin:       "lkslogin",
		AccordPassword: "accpass",
		LKSPassword:    "lkspass",
		IsActual:       false,
	}
	actualAccPass := "actualpassacc"
	actualLKSPass := "actuallkspass"

	getFn := s.dao.EXPECT().GetByID(mock.Anything, cr.ID).Return(cr, nil)
	defer getFn.Unset()
	checkFn := s.checker.EXPECT().Check(mock.Anything, cr.AccordLogin, actualAccPass, cr.LKSLogin, actualLKSPass).Return(nil)
	defer checkFn.Return(nil)
	saveFn := s.dao.EXPECT().Save(mock.Anything, cr).Return(nil)
	defer saveFn.Unset()

	err := s.serv.Update(context.Background(), cr.ID, actualAccPass, actualLKSPass)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), cr.AccordPassword, actualAccPass)
	assert.Equal(s.T(), cr.LKSPassword, actualLKSPass)
	assert.Equal(s.T(), cr.IsActual, true)
}

func (s *UpdateCredentialsSuite) TestCheckFailed() {
	cr := &entity.Credential{
		Model: gorm.Model{
			ID: 1,
		},
		AccordLogin:    "acclogin",
		LKSLogin:       "lkslogin",
		AccordPassword: "accpass",
		LKSPassword:    "lkspass",
		IsActual:       false,
	}
	actualAccPass := "actualpassacc"
	actualLKSPass := "actuallkspass"

	getFn := s.dao.EXPECT().GetByID(mock.Anything, cr.ID).Return(cr, nil)
	defer getFn.Unset()
	checkFn := s.checker.EXPECT().Check(mock.Anything, cr.AccordLogin, actualAccPass, cr.LKSLogin, actualLKSPass).Return(lks.ErrLKSAuth)
	defer checkFn.Return(nil)

	err := s.serv.Update(context.Background(), cr.ID, actualAccPass, actualLKSPass)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, lks.ErrLKSAuth)
}

func TestService_UpdateCredential(t *testing.T) {
	suite.Run(t, new(UpdateCredentialsSuite))
}
