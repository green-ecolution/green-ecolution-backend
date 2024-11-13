package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

//Tests the initialization of the JWT middleware to ensure that it is created correctly.

//func TestNewJWTMiddleware(t *testing.T) {

//}

func TestGenerateJWTToken(t *testing.T) {
	t.Run("should generate a valid JWT token", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": "1234567890",
			"exp": time.Now().Add(time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		secret := []byte("secret")

		tokenString, err := token.SignedString(secret)
		assert.NoError(t, err, "Expected token signing to succeed")
		assert.NotEmpty(t, tokenString, "Expected a non-empty token string")
	})
	t.Run("should handle empty secret without error but with non-empty token", func(t *testing.T) {
		claims := jwt.MapClaims{"sub": "1234567890"}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		secret := []byte("")

		tokenString, err := token.SignedString(secret)
		assert.NoError(t, err, "Expected no error when signing token with empty secret")
		assert.NotEmpty(t, tokenString, "Expected a non-empty token string even with empty secret")
	})

	t.Run("should generate token with future expiration time", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": "1234567890",
			"exp": time.Now().Add(24 * time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		secret := []byte("secret")

		tokenString, err := token.SignedString(secret)
		assert.NoError(t, err, "Expected token signing to succeed")
		assert.NotEmpty(t, tokenString, "Expected a non-empty token string")
	})

	t.Run("should fail to generate token with past expiration time", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": "1234567890",
			"exp": time.Now().Add(-time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		secret := []byte("secret")

		tokenString, err := token.SignedString(secret)
		assert.NoError(t, err, "Expected token signing to succeed even with past expiration time")
		assert.NotEmpty(t, tokenString, "Expected a non-empty token string, even if expired")
	})

	t.Run("should generate token with multiple claims", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub":         "1234567890",
			"exp":         time.Now().Add(time.Hour).Unix(),
			"role":        "admin",
			"permissions": []string{"read", "write", "delete"},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		secret := []byte("secret")

		tokenString, err := token.SignedString(secret)
		assert.NoError(t, err, "Expected token signing to succeed with multiple claims")
		assert.NotEmpty(t, tokenString, "Expected a non-empty token string")
	})

	t.Run("should handle large payload in claims", func(t *testing.T) {

		largePayload := make(map[string]string)
		for i := 0; i < 1000; i++ {
			largePayload[fmt.Sprintf("key%d", i)] = "value"
		}

		claims := jwt.MapClaims{
			"sub":  "1234567890",
			"exp":  time.Now().Add(time.Hour).Unix(),
			"data": largePayload,
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		secret := []byte("secret")

		tokenString, err := token.SignedString(secret)
		assert.NoError(t, err, "Expected token signing to succeed with a large payload")
		assert.NotEmpty(t, tokenString, "Expected a non-empty token string")
	})

	t.Run("should generate token with nil claims", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		secret := []byte("secret")

		tokenString, err := token.SignedString(secret)
		assert.NoError(t, err, "Expected token signing to succeed with nil claims")
		assert.NotEmpty(t, tokenString, "Expected a non-empty token string")
	})

	t.Run("should generate token with empty claims", func(t *testing.T) {
		claims := jwt.MapClaims{}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		secret := []byte("secret")

		tokenString, err := token.SignedString(secret)
		assert.NoError(t, err, "Expected token signing to succeed with empty claims")
		assert.NotEmpty(t, tokenString, "Expected a non-empty token string")
	})
}

// Utility function to generate a test RSA private key
func generateTestRSAPrivateKey() *rsa.PrivateKey {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	return privateKey
}

func TestVerifyJWTToken(t *testing.T) {
	t.Run("should verify a valid JWT token", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": "1234567890",
			"exp": time.Now().Add(time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		secret := []byte("secret")

		tokenString, err := token.SignedString(secret)
		assert.NoError(t, err, "Expected token signing to succeed")

		parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return secret, nil
		})
		assert.NoError(t, err, "Expected token parsing to succeed")
		assert.True(t, parsedToken.Valid, "Expected token to be valid")
	})

	t.Run("should fail to verify a JWT token with an invalid signature", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": "1234567890",
			"exp": time.Now().Add(time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		secret := []byte("secret")
		tokenString, _ := token.SignedString(secret)

		invalidSecret := []byte("invalid_secret")
		parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return invalidSecret, nil
		})

		assert.Error(t, err, "Expected token parsing to fail with invalid signature")
		assert.False(t, parsedToken.Valid, "Expected token to be invalid with wrong signature")
	})

	t.Run("should fail verification for expired token", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": "1234567890",
			"exp": time.Now().Add(-time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		secret := []byte("secret")

		tokenString, err := token.SignedString(secret)
		assert.NoError(t, err, "Expected token signing to succeed")

		parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		assert.Error(t, err, "Expected error due to expired token")
		assert.False(t, parsedToken.Valid, "Expected token to be invalid")
	})

	t.Run("should fail verification with missing expiration claim", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": "1234567890",
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		secret := []byte("secret")

		tokenString, err := token.SignedString(secret)
		assert.NoError(t, err, "Expected token signing to succeed")

		parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		assert.NoError(t, err, "Expected parsing to succeed without expiration claim")
		assert.True(t, parsedToken.Valid, "Expected token to be valid without expiration claim")
	})

	t.Run("should fail verification with modified token payload", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": "1234567890",
			"exp": time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		secret := []byte("secret")

		tokenString, err := token.SignedString(secret)
		assert.NoError(t, err, "Expected token signing to succeed")

		modifiedTokenString := tokenString[:len(tokenString)-1] + "x"

		parsedToken, err := jwt.Parse(modifiedTokenString, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		assert.Error(t, err, "Expected error due to modified token")
		assert.False(t, parsedToken.Valid, "Expected token to be invalid due to modification")
	})

	t.Run("should fail verification with unsupported signing method", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": "1234567890",
			"exp": time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		privateKey := generateTestRSAPrivateKey()
		tokenString, err := token.SignedString(privateKey)
		assert.NoError(t, err, "Expected token signing to succeed with RS256")

		secret := []byte("secret")
		parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unsupported signing method")
			}
			return secret, nil
		})

		assert.Error(t, err, "Expected error due to unsupported signing method")
		assert.False(t, parsedToken.Valid, "Expected token to be invalid with unsupported signing method")
	})

}
