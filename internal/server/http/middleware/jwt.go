package middleware

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"

	contribJwt "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	golangJwt "github.com/golang-jwt/jwt/v5"
	"github.com/green-ecolution/green-ecolution-backend/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils/enums"
	"github.com/pkg/errors"
)

func NewJWTMiddleware(cfg *config.IdentityAuthConfig, svc service.AuthService) fiber.Handler {
	base64Str := cfg.KeyCloak.RealmPublicKey
	publicKey, err := parsePublicKey(base64Str)
	if err != nil {
		return func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusInternalServerError).SendString("failed to parse public key")
		}
	}

	return contribJwt.New(contribJwt.Config{
		SigningKey: contribJwt.SigningKey{
			JWTAlg: contribJwt.RS256,
			Key:    publicKey,
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			return successHandler(c, svc)
		},
    ErrorHandler: func(c *fiber.Ctx, err error) error {
      return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
    },
	})
}

func parsePublicKey(base64Str string) (*rsa.PublicKey, error) {
	buf, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}

	parsedKey, err := x509.ParsePKIXPublicKey(buf)
	if err != nil {
		return nil, err
	}

	publicKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to parse public key")
	}

	return publicKey, nil
}

func successHandler(c *fiber.Ctx, svc service.AuthService) error {
	jwtToken := c.Locals("user").(*golangJwt.Token)
	claims := jwtToken.Claims.(golangJwt.MapClaims)

	ctx := c.Context()
	contextWithClaims := context.WithValue(ctx, enums.ContextKeyClaims, claims)
	c.SetUserContext(contextWithClaims)

	rptResult, err := svc.RetrospectToken(ctx, jwtToken.Raw)
	if err != nil {
		return err
	}

	if !*rptResult.Active {
		return c.Status(fiber.StatusUnauthorized).SendString("token is not active")
	}

	return c.Next()
}
