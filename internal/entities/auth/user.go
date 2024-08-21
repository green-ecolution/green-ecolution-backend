package auth

import (
	"net/url"
	"time"
)

type User struct {
	ID            string
	CreatedAt     time.Time
	Username      string `validate:"required,min=3,max=15"`
	FirstName     string `validate:"required,min=3,max=30"`
	LastName      string `validate:"required,min=3,max=30"`
	Email         string `validate:"required,email"`
	EmployeeID    string
	PhoneNumber   string
	EmailVerified bool
	Avatar        *url.URL
}

type RegisterUser struct {
	User     User
	Password string    `validate:"required"`
	Roles    *[]string 
}
