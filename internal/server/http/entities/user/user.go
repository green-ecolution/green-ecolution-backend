package user

import (
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/pagination"
)

type UserResponse struct {
	ID            string    `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	Username      string    `json:"username"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Email         string    `json:"email"`
	EmployeeID    string    `json:"employee_id"`
	PhoneNumber   string    `json:"phone_number"`
	EmailVerified bool      `json:"email_verified"`
	Avatar        string    `json:"avatar_url"`
} // @Name User

type UserListResponse struct {
	Data       []UserResponse        `json:"data"`
	Pagination pagination.Pagination `json:"pagination"`
} // @Name UserList

type UserRegisterRequest struct {
	Username    string    `json:"username"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	EmployeeID  string    `json:"employee_id,omitempty"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Password    string    `json:"password"`
	Roles       *[]string `json:"roles"`
	Avatar      string    `json:"avatar_url,omitempty"`
} // @Name UserRegister

type UserUpdateRequest struct {
	Username    string `json:"username,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Email       string `json:"email,omitempty"`
	EmployeeID  string `json:"employee_id,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Avatar      string `json:"avatar_url,omitempty"`
} // @Name UserUpdate

