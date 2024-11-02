package keycloak

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyCloakRepo_RetrospectToken(t *testing.T) {
	t.Run("should return valid token", func(t *testing.T) {
		// given
		ctx := context.Background()
		cfg := suite.IdentityConfig(t, ctx)
		k := NewKeycloakRepository(cfg)

		validToken, err := loginRestAPIClient(ctx, cfg.KeyCloak.BaseURL, cfg.KeyCloak.ClientID, cfg.KeyCloak.ClientSecret, cfg.KeyCloak.Realm)
		assert.NoError(t, err)

		// when
		got, err := k.RetrospectToken(ctx, validToken.AccessToken)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.NotNil(t, got.Active)
		assert.True(t, *got.Active)
	})

	t.Run("should return not active token when token is invalid", func(t *testing.T) {
		// given
		ctx := context.Background()
		cfg := suite.IdentityConfig(t, ctx)
		k := NewKeycloakRepository(cfg)

		// when
		got, err := k.RetrospectToken(ctx, "invalid-token")

		// then
		assert.NotNil(t, got)
		assert.NotNil(t, got.Active)
		assert.False(t, *got.Active)
		assert.NoError(t, err)
	})

	t.Run("should return not active token when token is empty", func(t *testing.T) {
		// given
		ctx := context.Background()
		cfg := suite.IdentityConfig(t, ctx)
		k := NewKeycloakRepository(cfg)

		// when
		got, err := k.RetrospectToken(ctx, "")

		// then
		assert.NotNil(t, got)
		assert.NotNil(t, got.Active)
		assert.False(t, *got.Active)
		assert.NoError(t, err)
	})
}



func TestKeyCloakRepo_GetAccessTokenFromClientCode(t *testing.T) {
}

func TestKeyCloakRepo_RefreshToken(t *testing.T) {
	t.Run("should return valid token", func(t *testing.T) {
	})

	t.Run("should return error when refresh token is invalid", func(t *testing.T) {
		// given
		ctx := context.Background()
		cfg := suite.IdentityConfig(t, ctx)
		k := NewKeycloakRepository(cfg)

		// when
		got, err := k.RefreshToken(ctx, "invalid-refresh-token")

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}
