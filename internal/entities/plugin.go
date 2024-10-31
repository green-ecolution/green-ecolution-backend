package entities

import (
	"net/url"
	"time"
)

type Plugin struct {
	Slug          string
	Name          string
	Path          url.URL
	LastHeartbeat time.Time
	Version       string
	Description   string
}
