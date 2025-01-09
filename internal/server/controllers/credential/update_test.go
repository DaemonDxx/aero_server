package ctrl_credential

import (
	"bytes"
	"encoding/json"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/api/lks"
	"github.com/daemondxx/lks_back/internal/server/middleware"
	ctrl_credential_mock "github.com/daemondxx/lks_back/mocks/server/controllers/credential"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type CredUpdateSuite struct {
	suite.Suite
	acc       *entity.Account
	serv      *ctrl_credential_mock.MockCredentialService
	ctrl      *controller
	w         *httptest.ResponseRecorder
	c         *gin.Context
	errMiddle gin.HandlerFunc
}

func (s *CredUpdateSuite) SetupSuite() {
	s.acc = &entity.Account{
		Model: gorm.Model{
			ID: 1,
		},
		TelegramID: 123456,
	}
	var crID uint = 1
	s.acc.CredentialID = &crID
	s.serv = &ctrl_credential_mock.MockCredentialService{}
	s.ctrl = &controller{
		crServ: s.serv,
	}
	s.errMiddle = middleware.ErrorMiddleware()
}

func (s *CredUpdateSuite) BeforeTest(suiteName, testName string) {
	s.w = httptest.NewRecorder()
	c, _ := gin.CreateTestContext(s.w)
	middleware.ProvideAccountInfo(c, s.acc)
	s.c = c
}

func TestController_UpdateCredential(t *testing.T) {
	suite.Run(t, new(CredUpdateSuite))
}

func (s *CredUpdateSuite) TestSuccessUpdate() {
	body := UpdateCredentialBody{
		AccordPassword: "testpassword",
		LKSPassword:    "testpassword",
	}
	bodyStr, _ := json.Marshal(body)
	var crID uint = 1

	createFn := s.serv.EXPECT().Update(
		mock.Anything,
		crID,
		body.AccordPassword,
		body.LKSPassword,
	).Return(nil)
	defer createFn.Unset()

	s.c.Params = append(s.c.Params, gin.Param{
		Key:   "id",
		Value: strconv.FormatUint(uint64(crID), 10),
	})
	s.c.Request = httptest.NewRequest(
		"PATCH",
		"/"+strconv.FormatUint(uint64(crID), 10),
		bytes.NewBuffer(bodyStr))
	s.ctrl.UpdateCredential(s.c)
	s.errMiddle(s.c)

	require.Equal(s.T(), http.StatusOK, s.w.Code)
}

func (s *CredUpdateSuite) TestUpdateAnotherCred() {
	body := UpdateCredentialBody{
		AccordPassword: "testpassword",
		LKSPassword:    "testpassword",
	}
	bodyStr, _ := json.Marshal(body)
	var crID uint = 2

	s.c.Params = append(s.c.Params, gin.Param{
		Key:   "id",
		Value: strconv.FormatUint(uint64(crID), 10),
	})
	s.c.Request = httptest.NewRequest(
		"PATCH",
		"/"+strconv.FormatUint(uint64(crID), 10),
		bytes.NewBuffer(bodyStr))

	s.ctrl.UpdateCredential(s.c)
	s.errMiddle(s.c)

	require.Equal(s.T(), http.StatusForbidden, s.w.Code)
}

func (s *CredUpdateSuite) TestUpdateWithoutID() {
	body := UpdateCredentialBody{
		AccordPassword: "testpassword",
		LKSPassword:    "testpassword",
	}
	bodyStr, _ := json.Marshal(body)
	s.c.Request = httptest.NewRequest(
		"PATCH",
		"/",
		bytes.NewBuffer(bodyStr))

	s.ctrl.UpdateCredential(s.c)
	s.errMiddle(s.c)

	require.Equal(s.T(), http.StatusBadRequest, s.w.Code)
}

func (s *CredUpdateSuite) TestUpdateStringID() {
	body := UpdateCredentialBody{
		AccordPassword: "testpassword",
		LKSPassword:    "testpassword",
	}
	bodyStr, _ := json.Marshal(body)
	s.c.Params = append(s.c.Params, gin.Param{
		Key:   "id",
		Value: "stringID",
	})

	s.c.Request = httptest.NewRequest(
		"PATCH",
		"/",
		bytes.NewBuffer(bodyStr))

	s.ctrl.UpdateCredential(s.c)
	s.errMiddle(s.c)

	require.Equal(s.T(), http.StatusBadRequest, s.w.Code)
}

func (s *CredUpdateSuite) TestFailureAuth() {
	body := UpdateCredentialBody{
		AccordPassword: "testpassword",
		LKSPassword:    "testpassword",
	}
	bodyStr, _ := json.Marshal(body)
	var crID uint = 1

	createFn := s.serv.EXPECT().Update(
		mock.Anything,
		crID,
		body.AccordPassword,
		body.LKSPassword,
	).Return(lks.ErrLKSAuth)
	defer createFn.Unset()

	s.c.Params = append(s.c.Params, gin.Param{
		Key:   "id",
		Value: strconv.FormatUint(uint64(crID), 10),
	})
	s.c.Request = httptest.NewRequest(
		"PATCH",
		"/"+strconv.FormatUint(uint64(crID), 10),
		bytes.NewBuffer(bodyStr))
	s.ctrl.UpdateCredential(s.c)
	s.errMiddle(s.c)

	require.Equal(s.T(), http.StatusForbidden, s.w.Code)

	var res *struct {
		Message string `json:"message"`
		System  string `json:"system"`
	}
	err := json.Unmarshal(s.w.Body.Bytes(), &res)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "LKS", res.System)
}
