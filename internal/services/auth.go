// internal/services/auth.go
package services

import (
	"context"
	"errors"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type authService struct {
	provider    *oidc.Provider
	verifier    *oidc.IDTokenVerifier
	oauthConfig oauth2.Config
}

func NewAuthService(provider *oidc.Provider, verifier *oidc.IDTokenVerifier, oauthConfig oauth2.Config) AuthService {
	return &authService{
		provider:    provider,
		verifier:    verifier,
		oauthConfig: oauthConfig,
	}
}

func (s *authService) Authenticate(token string) (string, error) {
	if strings.TrimSpace(token) == "" {
		return "", errors.New("empty token")
	}

	// Verify the token
	idToken, err := s.verifier.Verify(context.Background(), token)
	if err != nil {
		return "", err
	}

	// Extract claims
	var claims struct {
		Email string `json:"email"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return "", err
	}

	return claims.Email, nil
}
