package token

import "github.com/golang-jwt/jwt/v5"

type Config struct {
	Secret []byte
}

type Claims struct {
	AccountID uint `json:"accountID"`
	jwt.RegisteredClaims
}

type Service struct {
	cfg Config
}

func NewTokenService(c Config) *Service {
	return &Service{
		cfg: c,
	}
}
