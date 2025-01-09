package ctrl_credential

import (
	"bytes"
	"encoding/json"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/api/lks"
	"github.com/daemondxx/lks_back/internal/server/middleware"
	service_credential "github.com/daemondxx/lks_back/internal/server/services/credential"
	ctrl_credential_mock "github.com/daemondxx/lks_back/mocks/server/controllers/credential"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type CredCreateSuite struct {
	suite.Suite
	acc  *entity.Account
	serv *ctrl_credential_mock.MockCredentialService
	ctrl *controller
	w    *httptest.ResponseRecorder
	c    *gin.Context
}

func (s *CredCreateSuite) SetupSuite() {
	s.serv = &ctrl_credential_mock.MockCredentialService{}
	s.ctrl = &controller{
		crServ: s.serv,
	}
}

func (s *CredCreateSuite) BeforeTest(suiteName, testName string) {
	s.w = httptest.NewRecorder()
	c, _ := gin.CreateTestContext(s.w)
	s.acc = &entity.Account{
		Model: gorm.Model{
			ID: 1,
		},
		TelegramID:   123456,
		CredentialID: nil,
	}
	middleware.ProvideAccountInfo(c, s.acc)
	s.c = c
}

func (s *CredCreateSuite) TestSuccessCreate() {
	body := CreateCredentialsBody{
		AccordLogin:    "testlogin",
		AccordPassword: "testpassword",
		LKSLogin:       "testlogin",
		LKSPassword:    "testpassword",
	}
	bodyStr, _ := json.Marshal(body)

	createFn := s.serv.EXPECT().Create(
		mock.Anything,
		s.acc,
		body.AccordLogin,
		body.AccordPassword,
		body.LKSLogin,
		body.LKSPassword,
	).Return(&entity.Credential{
		Model: gorm.Model{
			ID: 1,
		},
		AccordLogin:    "testlogin",
		AccordPassword: "testpassword",
		LKSLogin:       "testlogin",
		LKSPassword:    "testpassword",
	}, nil)
	defer createFn.Unset()

	s.c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(bodyStr))
	s.ctrl.CreateCredential(s.c)

	require.Equal(s.T(), http.StatusCreated, s.w.Code)

	var res *CreateCredentialsResponse
	err := json.Unmarshal(s.w.Body.Bytes(), &res)
	require.NoError(s.T(), err)
	require.Equal(s.T(), body.LKSLogin, res.LksLogin)
	require.Equal(s.T(), body.AccordLogin, res.AccordLogin)
}

func (s *CredCreateSuite) TestFailedAccCredentialIsConnected() {
	body := CreateCredentialsBody{
		AccordLogin:    "testlogin",
		AccordPassword: "testpassword",
		LKSLogin:       "testlogin",
		LKSPassword:    "testpassword",
	}
	bodyStr, _ := json.Marshal(body)

	var crID uint = 1
	s.acc.CredentialID = &crID
	middleware.ProvideAccountInfo(s.c, s.acc)

	createFn := s.serv.EXPECT().Create(
		mock.Anything,
		s.acc,
		body.AccordLogin,
		body.AccordPassword,
		body.LKSLogin,
		body.LKSPassword,
	).Return(nil, service_credential.ErrCredentialIsConnected)
	defer createFn.Unset()

	s.c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(bodyStr))
	s.ctrl.CreateCredential(s.c)

	require.Equal(s.T(), http.StatusForbidden, s.w.Code)

	var res *struct {
		Message string `json:"message"`
	}
	err := json.Unmarshal(s.w.Body.Bytes(), &res)
	require.NoError(s.T(), err)
	require.Contains(s.T(), res.Message, "credential for this account is connected")
}

func (s *CredCreateSuite) TestFailedAuth() {
	body := CreateCredentialsBody{
		AccordLogin:    "testlogin",
		AccordPassword: "testpassword",
		LKSLogin:       "testlogin",
		LKSPassword:    "testpassword",
	}
	bodyStr, _ := json.Marshal(body)

	createFn := s.serv.EXPECT().Create(
		mock.Anything,
		s.acc,
		body.AccordLogin,
		body.AccordPassword,
		body.LKSLogin,
		body.LKSPassword,
	).Return(nil, lks.ErrLKSAuth)
	defer createFn.Unset()

	s.c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(bodyStr))
	s.ctrl.CreateCredential(s.c)

	require.Equal(s.T(), http.StatusForbidden, s.w.Code)

	var res *struct {
		Message string `json:"message"`
		System  string `json:"system"`
	}
	err := json.Unmarshal(s.w.Body.Bytes(), &res)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "LKS", res.System)
}

func TestController_CreateCredential(t *testing.T) {
	suite.Run(t, new(CredCreateSuite))
}
