package entities

import (
	"net/url"
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	UserStatusAvailable UserStatus = "available"
	UserStatusAbsent    UserStatus = "absent"
	UserStatusUnknown   UserStatus = "unknown"
)

type User struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	Username       string `validate:"required,min=3,max=15"`
	FirstName      string `validate:"required,min=3,max=30"`
	LastName       string `validate:"required,min=3,max=30"`
	Email          string `validate:"required,email"`
	EmployeeID     string
	PhoneNumber    string
	EmailVerified  bool
	Roles          []Role
	Avatar         *url.URL
	DrivingLicense DrivingLicense
	Status         UserStatus
}

type RegisterUser struct {
	User     User
	Password string `validate:"required"`
	Roles    []string
}

func ParseUserStatus(status string) UserStatus {
	switch status {
	case string(UserStatusAvailable):
		return UserStatus(UserStatusAvailable)
	case string(UserStatusAbsent):
		return UserStatus(UserStatusAbsent)
	default:
		return UserStatus(UserStatusUnknown)
	}
}
