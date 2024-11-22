package entities

import (
	"net/url"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	Username      string `validate:"required,min=3,max=15"`
	FirstName     string `validate:"required,min=3,max=30"`
	LastName      string `validate:"required,min=3,max=30"`
	Email         string `validate:"required,email"`
	EmployeeID    string
	PhoneNumber   string
	EmailVerified bool
	Avatar        *url.URL
	DriverLicense DriverLicense
}

type RegisterUser struct {
	User     User
	Password string `validate:"required"`
	Roles    []string
}
