package service

import (
	"crypto/sha256"
	"errors"

	"github.com/dwikie/auth-manager/internal/domain"
	"github.com/o1egl/paseto"
)

type ISecretProvider interface {
	GetAccessTokenSecret() string
	GetRefreshTokenSecret() string
}

type TokenService struct {
	secretProvider ISecretProvider
}

func NewTokenService(secretProvider ISecretProvider) *TokenService {
	return &TokenService{
		secretProvider: secretProvider,
	}
}

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

func (s *TokenService) GenerateAccessToken(claims *domain.TokenClaims, footer any) (string, error) {
	return s.generate(AccessToken, claims, footer)
}

func (s *TokenService) GenerateRefreshToken(claims *domain.TokenClaims, footer any) (string, error) {
	return s.generate(RefreshToken, claims, footer)
}

func (s *TokenService) ParseAccessToken(token string) (claims *domain.TokenClaims, footer any, err error) {
	return s.parse(AccessToken, token)
}

func (s *TokenService) ParseRefreshToken(token string) (claims *domain.TokenClaims, footer any, err error) {
	return s.parse(RefreshToken, token)
}

func (s *TokenService) generate(tokenType TokenType, claims *domain.TokenClaims, footer any) (string, error) {
	pasetoV2 := paseto.NewV2()

	symmetricKey, err := s.getSymmetricKey(tokenType)
	if err != nil {
		return "", err
	}

	token, err := pasetoV2.Encrypt(symmetricKey, claims, footer)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *TokenService) parse(tokenType TokenType, token string) (claims *domain.TokenClaims, footer any, err error) {
	errs := []error{}
	pasetoV2 := paseto.NewV2()

	symmetricKey, err := s.getSymmetricKey(tokenType)
	if err != nil {
		return nil, nil, errors.New("failed to get symmetric key: " + err.Error())
	}

	var parsedClaims domain.TokenClaims
	err = pasetoV2.Decrypt(token, symmetricKey, &parsedClaims, &footer)
	if err != nil {
		errs = append(errs, errors.New("failed to decrypt token: "+err.Error()))
		return nil, nil, errors.Join(errs...)
	}

	valid, err := parsedClaims.Verify()
	// fmt.Println("Parsed Claims:", valid, errs)
	if !valid {
		errs = append(errs, errors.New("token claims verification failed: "+err.Error()))
		return nil, nil, errors.Join(errs...)
	}

	return &parsedClaims, footer, nil
}

func (s *TokenService) getSymmetricKey(tokenType TokenType) ([]byte, error) {
	secretString := ""

	switch tokenType {
	case AccessToken:
		secretString = s.secretProvider.GetAccessTokenSecret()
	case RefreshToken:
		secretString = s.secretProvider.GetRefreshTokenSecret()
	default:
		return nil, errors.New("invalid token type")
	}

	hasher := sha256.New()
	hasher.Write([]byte(secretString))
	return hasher.Sum(nil), nil
}
