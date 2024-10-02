package entities

import (
	"net/url"
	"time"
)

type Plugin struct {
	Name          string
	Path          url.URL
	LastHeartbeat time.Time
  Version       string
  Description   string
}
