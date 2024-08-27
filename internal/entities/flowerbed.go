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
	Region         string
	Address        string
	Sensor         *Sensor
	Images         []*Image
	Latitude       float64
	Longitude      float64
}

type PlantingAreaWithImages struct {
	Flowerbed
	Image []*Image
}
