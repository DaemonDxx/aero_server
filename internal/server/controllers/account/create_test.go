package ctrl_account

import (
	"bytes"
	"encoding/json"
	"github.com/daemondxx/lks_back/internal/server/middleware"
	service_account "github.com/daemondxx/lks_back/internal/server/services/account"
	"github.com/daemondxx/lks_back/internal/services"
	account_mock "github.com/daemondxx/lks_back/mocks/server/controllers/account"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AccountCreateControllerSuite struct {
	suite.Suite
	accServ *account_mock.MockAccountService
	ctrl    *controller
	r       *gin.Engine
}

func (s *AccountCreateControllerSuite) SetupSuite() {
	s.accServ = &account_mock.MockAccountService{}
	s.ctrl = NewAccountController(s.accServ)
	s.r = gin.Default()
	s.r.Use(middleware.ErrorMiddleware())
	s.r.POST("/", s.ctrl.CreateAccount)
}

func (s *AccountCreateControllerSuite) BeforeTest(suiteName, testName string) {
}

func (s *AccountCreateControllerSuite) TestSuccess() {
	var tgID uint64 = 132
	token := "test.token"
	body, _ := json.Marshal(CreateAccountBody{TelegramID: tgID})

	createFn := s.accServ.EXPECT().CreateAccount(mock.Anything, tgID).Return(token, nil)
	defer createFn.Unset()

	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	s.r.ServeHTTP(w, req)

	require.Equal(s.T(), http.StatusCreated, w.Code)

	var res *CreateAccountResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(s.T(), err)
	require.Equal(s.T(), token, res.AccessToken)
}

func (s *AccountCreateControllerSuite) TestErrValidBody() {
	body := "{ \"telegramID\": \"test\" }"

	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	s.r.ServeHTTP(w, req)

	require.Equal(s.T(), http.StatusBadRequest, w.Code)
	require.Contains(s.T(), w.Body.String(), "parse json body error")

	body = "{}"
	req, _ = http.NewRequest("POST", "/", bytes.NewBufferString(body))
	s.r.ServeHTTP(w, req)
	require.Equal(s.T(), http.StatusBadRequest, w.Code)
	require.Contains(s.T(), w.Body.String(), "validation body error")
}

func (s *AccountCreateControllerSuite) TestAccIsExist() {
	var tgID uint64 = 132

	body, _ := json.Marshal(CreateAccountBody{TelegramID: tgID})

	createFn := s.accServ.EXPECT().CreateAccount(mock.Anything, tgID).Return("", &services.ErrServ{
		Service: "test",
		Message: "acc is exists",
		Err:     service_account.ErrAccountExists,
	})

	defer createFn.Unset()

	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	s.r.ServeHTTP(w, req)

	require.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func TestController_CreateAccount(t *testing.T) {
	suite.Run(t, new(AccountCreateControllerSuite))
}
