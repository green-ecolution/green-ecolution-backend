package keycloak

import (
	"context"
	"errors"

	"github.com/Nerzal/gocloak/v13"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrTokenExpired     = errors.New("token expired")
	ErrTokenNotActive   = errors.New("token not active")
	ErrTokenInvalidType = errors.New("token invalid type")
)

func (r *KeycloakRepository) RetrospectToken(ctx context.Context, token string) (*entities.IntroSpectTokenResult, error) {
	client := gocloak.NewClient(r.cfg.OidcProvider.BaseURL)

	rptResult, err := client.RetrospectToken(ctx, token, r.cfg.OidcProvider.Frontend.ClientID, r.cfg.OidcProvider.Frontend.ClientSecret, r.cfg.OidcProvider.DomainName)
	if err != nil {
		return nil, err
	}

	return &entities.IntroSpectTokenResult{
		Active:   rptResult.Active,
		Exp:      rptResult.Exp,
		AuthTime: rptResult.AuthTime,
		Type:     rptResult.Type,
	}, nil
}

func (r *KeycloakRepository) GetAccessTokenFromClientCode(ctx context.Context, code, redirectURL string) (*entities.ClientToken, error) {
	client := gocloak.NewClient(r.cfg.OidcProvider.BaseURL)

	tokenOptions := gocloak.TokenOptions{
		ClientID:     &r.cfg.OidcProvider.Frontend.ClientID,
		ClientSecret: &r.cfg.OidcProvider.Frontend.ClientSecret,
		Code:         &code,
		GrantType:    gocloak.StringP("authorization_code"),
		RedirectURI:  &redirectURL,
	}

	token, err := client.GetToken(ctx, r.cfg.OidcProvider.DomainName, tokenOptions)
	if err != nil {
		return nil, err
	}

	return &entities.ClientToken{
		AccessToken:      token.AccessToken,
		RefreshToken:     token.RefreshToken,
		ExpiresIn:        token.ExpiresIn,
		RefreshExpiresIn: token.RefreshExpiresIn,
		TokenType:        token.TokenType,
		NotBeforePolicy:  token.NotBeforePolicy,
		SessionState:     token.SessionState,
		Scope:            token.Scope,
	}, nil
}

func (r *KeycloakRepository) RefreshToken(ctx context.Context, refreshToken string) (*entities.ClientToken, error) {
	client := gocloak.NewClient(r.cfg.OidcProvider.BaseURL)
	token, err := client.RefreshToken(ctx, refreshToken, r.cfg.OidcProvider.Frontend.ClientID, r.cfg.OidcProvider.Frontend.ClientSecret, r.cfg.OidcProvider.DomainName)
	if err != nil {
		return nil, err
	}

	return &entities.ClientToken{
		AccessToken:      token.AccessToken,
		RefreshToken:     token.RefreshToken,
		ExpiresIn:        token.ExpiresIn,
		RefreshExpiresIn: token.RefreshExpiresIn,
		TokenType:        token.TokenType,
		NotBeforePolicy:  token.NotBeforePolicy,
		SessionState:     token.SessionState,
		Scope:            token.Scope,
	}, nil
}

func (r *KeycloakRepository) GetAccessTokenFromClientCredentials(ctx context.Context, clientID, clientSecret string) (*entities.ClientToken, error) {
	client := gocloak.NewClient(r.cfg.OidcProvider.BaseURL)

	tokenOptions := gocloak.TokenOptions{
		ClientID:     &clientID,
		ClientSecret: &clientSecret,
		GrantType:    gocloak.StringP("client_credentials"),
	}

	token, err := client.GetToken(ctx, r.cfg.OidcProvider.DomainName, tokenOptions)
	if err != nil {
		return nil, err
	}

	return &entities.ClientToken{
		AccessToken:      token.AccessToken,
		RefreshToken:     token.RefreshToken,
		ExpiresIn:        token.ExpiresIn,
		RefreshExpiresIn: token.RefreshExpiresIn,
		TokenType:        token.TokenType,
		NotBeforePolicy:  token.NotBeforePolicy,
		SessionState:     token.SessionState,
		Scope:            token.Scope,
	}, nil
}
