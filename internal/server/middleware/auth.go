package middleware

import (
	"context"
	"github.com/daemondxx/lks_back/entity"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

const injectToken = "ACC_TOKEN"

type TokenService interface {
	Parse(token string) (uint, error)
}

type AccountService interface {
	GetByID(ctx context.Context, id uint) (*entity.Account, error)
}

type AuthMiddleware struct {
	tServ   TokenService
	accServ AccountService
}

func NewAuthMiddleware(t TokenService, acc AccountService) *AuthMiddleware {
	return &AuthMiddleware{tServ: t, accServ: acc}
}

func (a *AuthMiddleware) Handler(c *gin.Context) {
	h := c.GetHeader("Authorization")
	if h == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "header \"Authorization\" not found",
		})
		return
	}

	_, t, ok := strings.Cut(h, "Bearer ")
	if !ok || t == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "header \"Authorization\" have not token",
		})
		return
	}

	accID, err := a.tServ.Parse(t)

	if err != nil {
		defer c.Abort()

		var msg string
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			msg = "token is malformed"
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			msg = "token signature is invalid"
		case errors.Is(err, jwt.ErrTokenExpired):
			msg = "token is expired"
		default:
			c.Error(errors.Wrap(err, "parse token error"))
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": msg,
		})
		return
	}

	acc, err := a.accServ.GetByID(c, accID)
	if err != nil {
		c.Error(errors.Wrap(err, "get account error"))
		return
	}

	ProvideAccountInfo(c, acc)
	c.Next()
}

func InjectAccountInfo(c *gin.Context) (*entity.Account, error) {
	v, ok := c.Get(injectToken)
	if !ok {
		return nil, errors.New("account info not provided")
	}

	acc, ok := v.(*entity.Account)
	if !ok {
		return nil, errors.New("account is not type entity.account")
	}

	return acc, nil
}

func ProvideAccountInfo(c *gin.Context, acc *entity.Account) {
	c.Set(injectToken, acc)
}
