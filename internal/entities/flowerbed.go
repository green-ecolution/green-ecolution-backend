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
	Archived       bool
	Latitude       float64
	Longitude      float64
}

type CreateFlowerbed struct {
	Size           float64
	Description    string
	NumberOfPlants int32
	MoistureLevel  float64
	Region         string
	Address        string
	Archived       bool
	Latitude       float64
	Longitude      float64
	SensorID       *int32
	ImageIDs       []int32
}

type UpdateFlowerbed struct {
  ID             int32
	Size           *float64
	Description    *string
	NumberOfPlants *int32
	MoistureLevel  *float64
	Region         *string
	Address        *string
	Archived       *bool
	Latitude       *float64
	Longitude      *float64
	SensorID       *int32
	ImageIDs       []int32
}

