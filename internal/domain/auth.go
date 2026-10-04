package domain

import (
	"errors"
	"time"
)

type TokenClaims struct {
	JTI        string         `json:"jti"`
	Issuer     string         `json:"iss"`
	Audience   string         `json:"aud"`
	Subject    string         `json:"sub"`
	Expiration time.Time      `json:"exp"`
	NotBefore  time.Time      `json:"nbf"`
	IssuedAt   time.Time      `json:"iat"`
	Extensions map[string]any `json:"ext,omitempty"`
}

func (t *TokenClaims) Verify() (valid bool, err error) {
	errs := []error{}
	if t.JTI == "" {
		errs = append(errs, errors.New("jti is required"))
	}
	if t.Issuer == "" {
		errs = append(errs, errors.New("iss is required"))
	}
	if t.Audience == "" {
		errs = append(errs, errors.New("aud is required"))
	}
	if t.Subject == "" {
		errs = append(errs, errors.New("sub is required"))
	}
	if t.Expiration.IsZero() {
		errs = append(errs, errors.New("exp is required"))
	}
	if t.NotBefore.IsZero() {
		errs = append(errs, errors.New("nbf is required"))
	}
	if t.IssuedAt.IsZero() {
		errs = append(errs, errors.New("iat is required"))
	}
	if time.Now().After(t.Expiration) {
		errs = append(errs, errors.New("token is expired"))
	}
	if time.Now().Before(t.NotBefore) {
		errs = append(errs, errors.New("token is not valid yet"))
	}

	return len(errs) == 0, errors.Join(errs...)
}
