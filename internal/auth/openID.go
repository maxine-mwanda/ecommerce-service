package auth

import (
	"context"
	"errors"
	"strings"

	"ecommerce-service/internal/config"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type OpenIDService struct {
	provider    *oidc.Provider
	verifier    *oidc.IDTokenVerifier
	oauthConfig oauth2.Config
}

/*func NewOpenIDService(cfg oauth2.Config, issuerURL string) *OpenIDService {
	provider, err := oidc.NewProvider(context.Background(), issuerURL)
	if err != nil {
		panic(err)
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.ClientID})

	oauthConfig := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  cfg.RedirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &OpenIDService{
		provider:    provider,
		verifier:    verifier,
		oauthConfig: oauthConfig,
	}
}*/

func NewOpenIDService(cfg config.AuthConfig) *OpenIDService {
	provider, err := oidc.NewProvider(context.Background(), cfg.IssuerURL)
	if err != nil {
		panic(err)
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.ClientID})

	oauthConfig := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &OpenIDService{
		provider:    provider,
		verifier:    verifier,
		oauthConfig: oauthConfig,
	}
}

func (s *OpenIDService) Authenticate(tokenString string) (string, error) {
	if strings.TrimSpace(tokenString) == "" {
		return "", errors.New("empty token")
	}

	token, err := s.verifier.Verify(context.Background(), tokenString)
	if err != nil {
		return "", err
	}

	var claims struct {
		Email string `json:"email"`
	}
	if err := token.Claims(&claims); err != nil {
		return "", err
	}

	return claims.Email, nil
}
