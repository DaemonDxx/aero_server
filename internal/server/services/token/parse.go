package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

func (s *Service) Parse(t string) (uint, error) {
	token, err := jwt.ParseWithClaims(t, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.cfg.Secret, nil
	})

	if err != nil {
		return 0, err
	}

	c, ok := token.Claims.(*Claims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	return c.AccountID, nil
}
