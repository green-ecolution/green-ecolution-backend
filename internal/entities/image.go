package entities

import "time"

type Image struct {
	ID        int32
	CreatedAt time.Time
	UpdatedAt time.Time
	URL       string
	Filename  *string
	MimeType  *string
}
