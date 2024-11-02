package keycloak

import (
	"context"
	"fmt"
	"testing"

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
