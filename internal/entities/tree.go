package entities

import (
	"time"
)

type Tree struct {
	ID             int32
	CreatedAt      time.Time
	UpdatedAt      time.Time
	TreeCluster    *TreeCluster
	Sensor         *Sensor
	Images         []*Image
	Readonly       bool
	PlantingYear   int32
	Species        string
	Number         string
	Latitude       float64
	Longitude      float64
	WateringStatus WateringStatus
}
type TreeCreate struct {
	TreeClusterID *int32
	Readonly      bool
	PlantingYear  int32
	Species       string
	Number        string
	Latitude      float64
	Longitude     float64
}

type TreeUpdate struct {
	TreeClusterID *int32  `validate:"omitempty,gt=0"`
	PlantingYear  int32   `validate:"omitempty,gt=0"`
	Species       string  `validate:"omitempty"`
	Number        string  `validate:"omitempty"`
	Latitude      float64 `validate:"omitempty,min=-90,max=90"`
	Longitude     float64 `validate:"omitempty,min=-180,max=180"`
}
