package token

import (
	"github.com/daemondxx/lks_back/entity"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func (s *Service) Create(acc *entity.Account) string {
	claims := &Claims{
		AccountID: acc.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(365 * 24 * 60 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, _ := token.SignedString(s.cfg.Secret)
	return t
}
