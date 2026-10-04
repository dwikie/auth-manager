package middleware

import (
	"github.com/dwikie/auth-manager/internal/domain"
)

type ITokenService interface {
	ParseAccessToken(token string) (claims *domain.TokenClaims, footer any, err error)
	ParseRefreshToken(token string) (claims *domain.TokenClaims, footer any, err error)
}

type Middleware struct {
	tokenService ITokenService
}

func RegisterMiddleware(tokenService ITokenService) *Middleware {
	return &Middleware{
		tokenService: tokenService,
	}
}
