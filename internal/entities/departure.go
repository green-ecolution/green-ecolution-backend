package entities

import "time"

type Departure struct {
	ID             int32
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Name           string
	Description    string
	Latitude       *float64
	Longitude      *float64
}
