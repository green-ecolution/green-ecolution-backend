package entities

import (
	"time"
)

type Flowerbed struct {
	ID             int32
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Size           float64
	Description    string
	NumberOfPlants int32
	MoistureLevel  float64
	Region         *Region
	Address        string
	Sensor         *Sensor
	Images         []*Image
	Archived       bool
	Latitude       *float64
	Longitude      *float64
}
