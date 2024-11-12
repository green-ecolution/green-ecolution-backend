package keycloak

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/google/uuid"
	"github.com/green-ecolution/green-ecolution-backend/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/stretchr/testify/assert"
)

var testUser = testUserToCreateFunc()

func TestKeyCloakUserRepo_Create(t *testing.T) {
	type fields struct {
		cfg *config.IdentityAuthConfig
	}
	type args struct {
		ctx      context.Context
		user     *entities.User
		password string
		roles    []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entities.User
		wantErr bool
	}{
		{
			name: "should create user successfully",
			fields: fields{
				cfg: suite.IdentityConfig(t, context.Background()),
			},
			args: args{
				ctx:      context.Background(),
				user:     testUser[0],
				password: "test",
				roles:    []string{},
			},
			wantErr: false,
		},
		{
			name: "should return error when create with same email",
			fields: fields{
				cfg: suite.IdentityConfig(t, context.Background()),
			},
			args: args{
				ctx:      context.Background(),
				user:     testUser[0],
				password: "test",
				roles:    []string{},
			},
			wantErr: true,
		},
		{
			name: "should return error when create with same username",
			fields: fields{
				cfg: suite.IdentityConfig(t, context.Background()),
			},
			args: args{
				ctx: context.Background(),
				user: &entities.User{
					Username:    testUser[0].Username,
					FirstName:   testUser[1].FirstName,
					LastName:    testUser[1].LastName,
					Email:       testUser[1].Email,
					EmployeeID:  testUser[1].EmployeeID,
					PhoneNumber: testUser[1].PhoneNumber,
				},
				password: "test",
				roles:    []string{},
			},
			wantErr: true,
		},
		{
			name: "should create user successfully with role",
			fields: fields{
				cfg: suite.IdentityConfig(t, context.Background()),
			},
			args: args{
				ctx:      context.Background(),
				user:     testUser[2],
				password: "test",
				roles:    []string{"admin"},
			},
			wantErr: false,
		},
		{
			name: "should create user successfully when roles are nil",
			fields: fields{
				cfg: suite.IdentityConfig(t, context.Background()),
			},
			args: args{
				ctx:      context.Background(),
				user:     testUser[3],
				password: "test",
				roles:    nil,
			},
			wantErr: false,
		},
		{
			name: "should return error when user is nil",
			fields: fields{
				cfg: suite.IdentityConfig(t, context.Background()),
			},
			args: args{
				ctx:      context.Background(),
				user:     nil,
				password: "test",
				roles:    nil,
			},
			wantErr: true,
		},
		{
			name: "should return error when role not exist",
			fields: fields{
				cfg: suite.IdentityConfig(t, context.Background()),
			},
			args: args{
				ctx:      context.Background(),
				user:     testUser[4],
				password: "test",
				roles:    []string{"not-exist"},
			},
			wantErr: true,
		},
		{
			name: "should return error when failed to set password",
			fields: fields{
				cfg: suite.IdentityConfig(t, context.Background()),
			},
			args: args{
				ctx:      context.Background(),
				user:     testUser[5],
				password: "",
				roles:    nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UserRepository{
				cfg: tt.fields.cfg,
			}
			got, err := r.Create(tt.args.ctx, tt.args.user, tt.args.password, tt.args.roles)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.args.user.Username, got.Username)
				assert.Equal(t, tt.args.user.FirstName, got.FirstName)
				assert.Equal(t, tt.args.user.LastName, got.LastName)
				assert.Equal(t, tt.args.user.Email, got.Email)
				assert.Equal(t, tt.args.user.EmployeeID, got.EmployeeID)
				assert.Equal(t, tt.args.user.PhoneNumber, got.PhoneNumber)
				assert.NotNil(t, got.ID)
				assert.NotZero(t, got.ID)
				assert.NotZero(t, got.CreatedAt)
			}
		})
	}
}

func TestKeyCloakUserRepo_RemoveSession(t *testing.T) {
	t.Run("should remove session successfully", func(t *testing.T) {
		// given
		identityConfig := suite.IdentityConfig(t, context.Background())
		userRepo := NewUserRepository(identityConfig)
		user := &entities.User{
			Username:    "should-remove-session",
			FirstName:   "Toni",
			LastName:    "Tester",
			Email:       "should-remove-session@green-ecolution.de",
			EmployeeID:  "123456",
			PhoneNumber: "+49 123456",
		}
		ensureUserExists(t, user)
		token := loginUser(t, user)

		// when
		err := userRepo.RemoveSession(context.Background(), token.RefreshToken)

		// then
		assert.NoError(t, err)
	})

	t.Run("should return error when failed to remove session", func(t *testing.T) {
		// given
		identityConfig := &config.IdentityAuthConfig{}
		userRepo := NewUserRepository(identityConfig)

		// when
		err := userRepo.RemoveSession(context.Background(), "invalid-token")

		// then
		assert.Error(t, err)
	})
}

func TestKeyCloakUserRepo_KeyCloakUserToUser(t *testing.T) {
	t.Run("should convert keycloak user to user successfully", func(t *testing.T) {
		// given
		uuid, _ := uuid.NewRandom()
		user := &gocloak.User{
			ID:               gocloak.StringP(uuid.String()),
			CreatedTimestamp: gocloak.Int64P(123456),
			Username:         gocloak.StringP("test"),
			FirstName:        gocloak.StringP("Toni"),
			LastName:         gocloak.StringP("Tester"),
			Email:            gocloak.StringP("dev@green-ecolution.de"),
			Attributes: &map[string][]string{
				"phone_number": {"+49 123456"},
				"employee_id":  {"123456"},
			},
		}

		// when
		got, err := keyCloakUserToUser(user)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, uuid.String(), got.ID.String())
		assert.Equal(t, "test", got.Username)
		assert.Equal(t, "Toni", got.FirstName)
		assert.Equal(t, "Tester", got.LastName)
		assert.Equal(t, "dev@green-ecolution.de", got.Email)
		assert.Equal(t, "+49 123456", got.PhoneNumber)
		assert.Equal(t, "123456", got.EmployeeID)
	})

	t.Run("should return error when failed to parse user id", func(t *testing.T) {
		// given
		user := &gocloak.User{
			ID:               gocloak.StringP("invalid-id"),
			CreatedTimestamp: gocloak.Int64P(123456),
			Username:         gocloak.StringP("test"),
			FirstName:        gocloak.StringP("Toni"),
			LastName:         gocloak.StringP("Tester"),
			Email:            gocloak.StringP("dev@green-ecolution.de"),
			Attributes: &map[string][]string{
				"phone_number": {"+49 123456"},
				"employee_id":  {"123456"},
			},
		}

		// when
		got, err := keyCloakUserToUser(user)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestKeyCloakUserRepo_UserToKeyCloakUser(t *testing.T) {
	t.Run("should convert user to keycloak user successfully", func(t *testing.T) {
		// given
		uuid, _ := uuid.NewRandom()
		user := &entities.User{
			ID:          uuid,
			CreatedAt:   time.Unix(123456, 0),
			Username:    "test",
			FirstName:   "Toni",
			LastName:    "Tester",
			Email:       "dev@green-ecolution.de",
			PhoneNumber: "+49 123456",
			EmployeeID:  "123456",
		}

		// when
		got := userToKeyCloakUser(user)

		// then
		assert.NotNil(t, got)
		assert.NotNil(t, got.ID)
		assert.NotNil(t, got.Username)
		assert.NotNil(t, got.FirstName)
		assert.NotNil(t, got.LastName)
		assert.NotNil(t, got.Email)
		assert.NotNil(t, got.Attributes)

		assert.Equal(t, uuid.String(), *got.ID)
		assert.Equal(t, "test", *got.Username)
		assert.Equal(t, "Toni", *got.FirstName)
		assert.Equal(t, "Tester", *got.LastName)
		assert.Equal(t, "dev@green-ecolution.de", *got.Email)
		assert.Equal(t, "+49 123456", (*got.Attributes)["phone_number"][0])
		assert.Equal(t, "123456", (*got.Attributes)["employee_id"][0])
	})
}

func loginAdminAndGetToken(t testing.TB) *gocloak.JWT {
	t.Helper()
	identityConfig := suite.IdentityConfig(t, context.Background())
	client := gocloak.NewClient(identityConfig.KeyCloak.BaseURL)
	token, err := client.Login(context.Background(), identityConfig.KeyCloak.ClientID, identityConfig.KeyCloak.ClientSecret, identityConfig.KeyCloak.Realm, suite.User, suite.Password)
	if err != nil {
		t.Fatalf("loginAdminAndGetToken::failed to get token: %v", err)
	}
	return token
}

func loginUser(t testing.TB, user *entities.User) *gocloak.JWT {
	t.Helper()
	identityConfig := suite.IdentityConfig(t, context.Background())
	client := gocloak.NewClient(identityConfig.KeyCloak.BaseURL)
	token, err := client.Login(context.Background(), identityConfig.KeyCloak.Frontend.ClientID, identityConfig.KeyCloak.Frontend.ClientSecret, identityConfig.KeyCloak.Realm, user.Username, "test")
	if err != nil {
		t.Fatalf("loginUser::failed to get token: %v", err)
	}

	return token
}

func ensureUserExists(t testing.TB, user *entities.User) {
	t.Helper()
	identityConfig := suite.IdentityConfig(t, context.Background())
	client := gocloak.NewClient(identityConfig.KeyCloak.BaseURL)
	token, err := client.LoginClient(context.Background(), identityConfig.KeyCloak.ClientID, identityConfig.KeyCloak.ClientSecret, identityConfig.KeyCloak.Realm)

  userID, err := client.CreateUser(context.Background(), token.AccessToken, identityConfig.KeyCloak.Realm, *userToKeyCloakUser(user))
	if err != nil {
		t.Log("ensureUserExists::failed to create user. maybe user already exists. error: ", err)
	}

	if err = client.SetPassword(context.Background(), token.AccessToken, userID, identityConfig.KeyCloak.Realm, "test", false); err != nil {
    t.Fatalf("ensureUserExists::failed to set password: %v", err)
	}
}

func testUserToCreateFunc() []*entities.User {
	n := 20
	users := make([]*entities.User, n)
	for i := 0; i < n; i++ {
		users[i] = &entities.User{
			Username:    fmt.Sprintf("test%d", i),
			FirstName:   "Toni",
			LastName:    "Tester",
			Email:       fmt.Sprintf("test%d@green-ecolution.de", i),
			EmployeeID:  "123456",
			PhoneNumber: "+49 123456",
		}
	}
	return users
}
