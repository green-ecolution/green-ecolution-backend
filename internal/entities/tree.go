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
	Description    string
}

type TreeCreate struct {
	TreeClusterID *int32  `validate:"required"`
	Readonly      bool    `validate:"omitempty"`
	PlantingYear  int32   `validate:"required"`
	Species       string  `validate:"required"`
	Number        string  `validate:"required"`
	Latitude      float64 `validate:"required,max=90,min=-90"`
	Longitude     float64 `validate:"required,max=180,min=-180"`
	Description   string  `validate:"omitempty"`
}

type TreeUpdate struct {
	TreeClusterID *int32  `validate:"omitempty,gt=0"`
	PlantingYear  int32   `validate:"omitempty,gt=0"`
	Species       string  `validate:"omitempty"`
	Number        string  `validate:"omitempty"`
	Latitude      float64 `validate:"omitempty,min=-90,max=90"`
	Longitude     float64 `validate:"omitempty,min=-180,max=180"`
	Description   string  `validate:"omitempty"`
}

type TreeImport struct {
	Number       string  `validate:"required"`
	Species      string  `validate:"omitempty"`
	Latitude     float64 `validate:"required,max=90,min=-90"`
	Longitude    float64 `validate:"required,max=180,min=-180"`
	PlantingYear int32   `validate:"omitempty,gt=0"`
	Street       string  `validate:"required"`
	TreeID       int32   `validate:"omitempty"`
}
