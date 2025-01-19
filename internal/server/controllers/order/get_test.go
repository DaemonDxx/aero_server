package order

import (
	"encoding/json"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/server/middleware"
	ctrl_order_mock "github.com/daemondxx/lks_back/mocks/server/controllers/order"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type GetOrderSuite struct {
	suite.Suite
	acc    *entity.Account
	orders []entity.Order
	serv   *ctrl_order_mock.MockOrderService
	ctrl   *controller
	w      *httptest.ResponseRecorder
	c      *gin.Context
}

func (s *GetOrderSuite) SetupSuite() {
	s.serv = &ctrl_order_mock.MockOrderService{}
	s.orders = []entity.Order{{
		Model: gorm.Model{
			ID: 1,
		},
	},
	}
	s.ctrl = &controller{
		serv: s.serv,
	}
}

func (s *GetOrderSuite) BeforeTest(suiteName, testName string) {
	var crID uint = 4

	s.w = httptest.NewRecorder()
	c, _ := gin.CreateTestContext(s.w)

	s.acc = &entity.Account{
		Model: gorm.Model{
			ID: 1,
		},
		TelegramID:   123456,
		CredentialID: &crID,
	}
	middleware.ProvideAccountInfo(c, s.acc)
	s.c = c
}

func (s *GetOrderSuite) TestGetOrder() {
	getFn := s.serv.EXPECT().GetActualOrders(
		mock.Anything,
		mock.Anything,
	).Return(s.orders, nil)
	defer getFn.Unset()

	s.c.Request = httptest.NewRequest("GET", "/", nil)
	s.ctrl.Get(s.c)

	require.Nil(s.T(), s.c.Errors.Last())
	require.Equal(s.T(), http.StatusOK, s.w.Code)

	var res *GetOrdersResponse
	err := json.Unmarshal(s.w.Body.Bytes(), &res)
	require.NoError(s.T(), err)
	require.Equal(s.T(), len(s.orders), len(res.Orders))
}

func (s *GetOrderSuite) TestGetOrderWithEmptyCredential() {
	s.acc.CredentialID = nil
	getFn := s.serv.EXPECT().GetActualOrders(
		mock.Anything,
		mock.Anything,
	).Return(s.orders, nil)
	defer getFn.Unset()

	s.c.Request = httptest.NewRequest("GET", "/", nil)
	s.ctrl.Get(s.c)

	require.Error(s.T(), s.c.Errors.Last())
}

func TestController_GetOrder(t *testing.T) {
	suite.Run(t, new(GetOrderSuite))
}
