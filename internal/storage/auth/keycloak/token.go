package keycloak

import (
	"context"
	"errors"

	"github.com/Nerzal/gocloak/v13"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities/auth"
)

var (
  ErrInvalidToken = errors.New("invalid token")
  ErrTokenExpired = errors.New("token expired")
  ErrTokenNotActive = errors.New("token not active")
  ErrTokenInvalidType = errors.New("token invalid type")
) 

func (r *KeycloakRepository) RetrospectToken(ctx context.Context, token string) (*auth.IntroSpectTokenResult, error) {
	client := gocloak.NewClient(r.cfg.KeyCloak.BaseURL)

	rptResult, err := client.RetrospectToken(ctx, token, r.cfg.KeyCloak.Frontend.ClientID, r.cfg.KeyCloak.Frontend.ClientSecret, r.cfg.KeyCloak.Realm)
	if err != nil {
		return nil, err
	}

	return &auth.IntroSpectTokenResult{
		Active:   rptResult.Active,
		Exp:      rptResult.Exp,
		AuthTime: rptResult.AuthTime,
		Type:     rptResult.Type,
	}, nil
}

func (r *KeycloakRepository) GetAccessTokenFromClientCode(ctx context.Context, code, redirectURL string) (*auth.ClientToken, error) {
  client := gocloak.NewClient(r.cfg.KeyCloak.BaseURL)

  tokenOptions := gocloak.TokenOptions{
    ClientID:     &r.cfg.KeyCloak.Frontend.ClientID,
    ClientSecret: &r.cfg.KeyCloak.Frontend.ClientSecret,
    Code:         &code,
    GrantType:   gocloak.StringP("authorization_code"),
    RedirectURI: &redirectURL,
  }

  token, err := client.GetToken(ctx, r.cfg.KeyCloak.Realm, tokenOptions)
  if err != nil {
    return nil, err
  }

  return &auth.ClientToken{
    AccessToken:  token.AccessToken,
    RefreshToken: token.RefreshToken,
    ExpiresIn:    token.ExpiresIn,
    RefreshExpiresIn: token.RefreshExpiresIn,
    TokenType: token.TokenType,
    NotBeforePolicy: token.NotBeforePolicy,
    SessionState: token.SessionState,
    Scope: token.Scope,
  }, nil
}
