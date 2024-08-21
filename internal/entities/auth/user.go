package auth

import (
	"net/url"
)

type User struct {
  ID              string   
	Username        string   `validate:"required,min=3,max=15"`
	Password        string   `validate:"required"`
	FirstName       string   `validate:"required,min=3,max=30"`
	LastName        string   `validate:"required,min=3,max=30"`
	Email           string   `validate:"required,email"`
	PhoneNumber     string   `validate:"required,min=10,max=15"`
	EmployeeID      string   `validate:"required"`
	ProfileImageURL *url.URL `validate:"required,url"`
}

type RegisterUser struct {
	User
	Password string `validate:"required"`
	Role     string `validate:"required"`
}
