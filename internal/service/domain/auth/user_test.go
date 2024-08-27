package auth

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/green-ecolution/green-ecolution-backend/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	t.Run("should return result when success", func(t *testing.T) {
		// given
		identityConfig := &config.IdentityAuthConfig{}
		inputUser := &entities.User{
			Username:    "username",
			Email:       "mail@foo.com",
			PhoneNumber: "phoneNumber",
			FirstName:   "firstName",
			LastName:    "lastName",
			EmployeeID:  "employeeID",
		}
		input := &entities.RegisterUser{
			User:     *inputUser,
			Password: "password",
			Roles:    &[]string{"viewer"},
		}

		expected := &entities.User{
			ID:            uuid.MustParse("6be4c752-94df-4719-99b1-ce58253eaf75"),
			CreatedAt:     time.Now(),
			Username:      inputUser.Username,
			FirstName:     inputUser.FirstName,
			LastName:      inputUser.LastName,
			Email:         inputUser.Email,
			EmployeeID:    inputUser.EmployeeID,
			PhoneNumber:   inputUser.PhoneNumber,
			EmailVerified: false,
			Avatar:        nil,
		}

		repo := storageMock.NewMockAuthRepository(t)
		svc := NewAuthService(repo, identityConfig)

		// when
		repo.EXPECT().CreateUser(context.Background(), inputUser, input.Password, input.Roles).Return(expected, nil)
		resp, err := svc.Register(context.Background(), input)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expected, resp)
	})

	t.Run("should return error when failed to register user", func(t *testing.T) {
		// given
		identityConfig := &config.IdentityAuthConfig{}
		repo := storageMock.NewMockAuthRepository(t)
		svc := NewAuthService(repo, identityConfig)

		// when
		resp, err := svc.Register(context.Background(), &entities.RegisterUser{})

		// then
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("should return error when validation error", func(t *testing.T) {
		// given
		identityConfig := &config.IdentityAuthConfig{}
		repo := storageMock.NewMockAuthRepository(t)
		svc := NewAuthService(repo, identityConfig)

		// when
		resp, err := svc.Register(context.Background(), &entities.RegisterUser{})

		// then
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
