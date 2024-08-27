package auth

import (
	"time"
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

type RegisterUserRequest struct {
	User     UserResponse `json:"user"`
	Password string       `json:"password"`
	Roles    *[]string    `json:"roles"`
} // @Name RegisterUser
