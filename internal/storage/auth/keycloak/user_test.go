package keycloak

import (
	"context"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestKeyCloakUserRepo_Create(t *testing.T) {
	t.Run("should create user successfully", func(t *testing.T) {
		// given
		ctx := context.Background()
		cfg := suite.IdentityConfig(t, ctx)
		repo := NewUserRepository(cfg)

		password := "test"
		user := &entities.User{
			Username:    "test",
			FirstName:   "Toni",
			LastName:    "Tester",
			Email:       "test@green-ecolution.de",
			EmployeeID:  "123456",
			PhoneNumber: "123456",
		}

		// when
		createdUser, err := repo.Create(ctx, user, password, []string{})

		// then
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.Equal(t, user.Username, createdUser.Username)
		assert.Equal(t, user.FirstName, createdUser.FirstName)
		assert.Equal(t, user.LastName, createdUser.LastName)
		assert.Equal(t, user.Email, createdUser.Email)
		// assert.Equal(t, user.EmployeeID, createdUser.EmployeeID) // fixme: handle additional fields in repo
		// assert.Equal(t, user.PhoneNumber, createdUser.PhoneNumber) // fixme: handle additional fields in repo
		assert.NotNil(t, createdUser.ID)
		assert.NotZero(t, createdUser.ID)
		assert.NotZero(t, createdUser.CreatedAt)
	})

	t.Run("should return error when create with same email", func(t *testing.T) {
		// given
		ctx := context.Background()
		cfg := suite.IdentityConfig(t, ctx)
		repo := NewUserRepository(cfg)

		password := "test"
		user := &entities.User{
			Username:    "test",
			FirstName:   "Toni",
			LastName:    "Tester",
			Email:       "test@green-ecolution.de",
			EmployeeID:  "123456",
			PhoneNumber: "123456",
		}

		// when
		createdUser, err := repo.Create(ctx, user, password, []string{})

		// then
		assert.Error(t, err)
		assert.Nil(t, createdUser)
	})

	t.Run("should return error when create with same username", func(t *testing.T) {
		// given
		ctx := context.Background()
		cfg := suite.IdentityConfig(t, ctx)
		repo := NewUserRepository(cfg)

		password := "test"
		user := &entities.User{
			Username:    "test",
			FirstName:   "Toni",
			LastName:    "Tester",
			Email:       "test1@green-ecolution.de",
			EmployeeID:  "123456",
			PhoneNumber: "123456",
		}

		// when
		createdUser, err := repo.Create(ctx, user, password, []string{})

		// then
		assert.Error(t, err)
		assert.Nil(t, createdUser)
	})

	t.Run("should create user successfully with role", func(t *testing.T) {
		// given
		ctx := context.Background()
		cfg := suite.IdentityConfig(t, ctx)
		repo := NewUserRepository(cfg)

		password := "test"
		user := &entities.User{
			Username:    "test2",
			FirstName:   "Toni",
			LastName:    "Tester",
			Email:       "test2@green-ecolution.de",
			EmployeeID:  "123456",
			PhoneNumber: "123456",
		}

		// when
		createdUser, err := repo.Create(ctx, user, password, []string{"admin"})

		// then
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.Equal(t, user.Username, createdUser.Username)
		assert.Equal(t, user.FirstName, createdUser.FirstName)
		assert.Equal(t, user.LastName, createdUser.LastName)
		assert.Equal(t, user.Email, createdUser.Email)
		// assert.Equal(t, user.EmployeeID, createdUser.EmployeeID) // fixme: handle additional fields in repo
		// assert.Equal(t, user.PhoneNumber, createdUser.PhoneNumber) // fixme: handle additional fields in repo
		assert.NotNil(t, createdUser.ID)
		assert.NotZero(t, createdUser.ID)
		assert.NotZero(t, createdUser.CreatedAt)
	})

	t.Run("should create user successfully when roles are nil", func(t *testing.T) {
		// given
		ctx := context.Background()
		cfg := suite.IdentityConfig(t, ctx)
		repo := NewUserRepository(cfg)

		password := "test"
		user := &entities.User{
			Username:    "test3",
			FirstName:   "Toni",
			LastName:    "Tester",
			Email:       "test3@green-ecolution.de",
			EmployeeID:  "123456",
			PhoneNumber: "123456",
		}

		// when
		createdUser, err := repo.Create(ctx, user, password, nil)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.Equal(t, user.Username, createdUser.Username)
		assert.Equal(t, user.FirstName, createdUser.FirstName)
		assert.Equal(t, user.LastName, createdUser.LastName)
		assert.Equal(t, user.Email, createdUser.Email)
		// assert.Equal(t, user.EmployeeID, createdUser.EmployeeID) // fixme: handle additional fields in repo
		// assert.Equal(t, user.PhoneNumber, createdUser.PhoneNumber) // fixme: handle additional fields in repo
		assert.NotNil(t, createdUser.ID)
		assert.NotZero(t, createdUser.ID)
		assert.NotZero(t, createdUser.CreatedAt)
	})

	t.Run("should return error when user is nil", func(t *testing.T) {
		// given
		ctx := context.Background()
		cfg := suite.IdentityConfig(t, ctx)
		repo := NewUserRepository(cfg)

		// when
		createdUser, err := repo.Create(ctx, nil, "test", nil)

		// then
		assert.Error(t, err)
		assert.Nil(t, createdUser)
	})

	t.Run("should return error when role not exist", func(t *testing.T) {
		// given
		ctx := context.Background()
		cfg := suite.IdentityConfig(t, ctx)
		repo := NewUserRepository(cfg)

		password := "test"
		user := &entities.User{
			Username:    "test4",
			FirstName:   "Toni",
			LastName:    "Tester",
			Email:       "test4@green-ecolution.de",
			EmployeeID:  "123456",
			PhoneNumber: "123456",
		}

		// when
		createdUser, err := repo.Create(ctx, user, password, []string{"not-exist"})

		// then
		assert.Error(t, err)
		assert.Nil(t, createdUser)
	})
}
