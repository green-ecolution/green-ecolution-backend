package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities/auth"
	"github.com/pkg/errors"
)

func (r *KeycloakRepository) RetrospectToken(ctx context.Context, token string) (*auth.IntroSpectTokenResult, error) {
  client := gocloak.NewClient(r.cfg.KeyCloak.BaseURL)

  rptResult, err := client.RetrospectToken(ctx, token, r.cfg.KeyCloak.ClientID, r.cfg.KeyCloak.ClientSecret, r.cfg.KeyCloak.Realm)
  if err != nil {
    return nil, errors.Wrap(err, "failed to retrospect token")
  }

  return &auth.IntroSpectTokenResult{
    Active: rptResult.Active,
    Exp:    rptResult.Exp,
    AuthTime: rptResult.AuthTime,
    Type:   rptResult.Type,
  }, nil
}
