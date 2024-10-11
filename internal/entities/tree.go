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
	TreeClusterID *int32
	Readonly      bool
	PlantingYear  int32
	Species       string
	Number        string
	Latitude      float64
	Longitude     float64
}
