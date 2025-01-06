package middleware

import (
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/server/services/token"
	middleware_mock "github.com/daemondxx/lks_back/mocks/server/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AuthMiddlewareSuite struct {
	suite.Suite
	tServ   *token.Service
	accServ *middleware_mock.MockAccountService
	r       *gin.Engine
}

func (s *AuthMiddlewareSuite) SetupSuite() {
	s.tServ = token.NewTokenService(token.Config{Secret: []byte("test")})
	s.accServ = &middleware_mock.MockAccountService{}
	authM := NewAuthMiddleware(s.tServ, s.accServ)

	r := gin.Default()
	r.Use(authM.Handler)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{})
	})
	s.r = r
}

func (s *AuthMiddlewareSuite) TestSuccessAuth() {
	acc := &entity.Account{
		Model: gorm.Model{
			ID: 1,
		},
	}
	token := s.tServ.Create(acc)

	accFn := s.accServ.EXPECT().GetByID(mock.Anything, acc.ID).Return(acc, nil)
	defer accFn.Unset()

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	s.r.ServeHTTP(w, req)

	require.Equal(s.T(), http.StatusOK, w.Code)
}

func (s *AuthMiddlewareSuite) TestInvalidHeader() {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer ")

	w := httptest.NewRecorder()

	s.r.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)

	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "test.token")
	s.r.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)

	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer test,token")
	s.r.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
}

func (s *AuthMiddlewareSuite) TestInjectAccount() {
	acc := &entity.Account{
		Model: gorm.Model{
			ID: 1,
		},
	}
	token := s.tServ.Create(acc)

	accFn := s.accServ.EXPECT().GetByID(mock.Anything, acc.ID).Return(acc, nil)
	defer accFn.Unset()

	s.r.GET("/inject", func(c *gin.Context) {
		a, err := InjectAccountInfo(c)
		require.NoError(s.T(), err)
		require.Equal(s.T(), acc.ID, a.ID)
	})

	req, _ := http.NewRequest("GET", "/inject", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	s.r.ServeHTTP(w, req)

	require.Equal(s.T(), http.StatusOK, w.Code)
}

func TestAuthMiddleware_Handler(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareSuite))
}
