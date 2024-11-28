package entities

import (
	"net/url"
)

type Plugin struct {
	Slug        string  `validate:"required"`
	Name        string  `validate:"required"`
	Path        url.URL
	Version     string
	Description string
	Auth        struct {
		Username string
		Password string
	}
}
