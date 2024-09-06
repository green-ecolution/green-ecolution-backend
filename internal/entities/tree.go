package entities

import (
	"time"
)

type Tree struct {
	ID                  int32
	CreatedAt           time.Time
	UpdatedAt           time.Time
	TreeCluster         *TreeCluster
	Sensor              *Sensor
	Images              []*Image
	Age                 int32
	HeightAboveSeaLevel float64
	PlantingYear        int32
	Species             string
	Number              int32
	Latitude            float64
	Longitude           float64
}

type CreateTree struct {
	TreeClusterID       int32
	Species             string
	Age                 int32
	HeightAboveSeaLevel float64
	PlantingYear        int32
	Latitude            float64
	Longitude           float64
	SensorID            *int32
	ImageIDs            []*int32
}

type UpdateTree struct {
	ID                  int32
	Species             *string
	Age                 *int32
	HeightAboveSeaLevel *float64
	SensorID            *int32
	ImageIDs            []*int32
	PlantingYear        *int32
	Latitude            *float64
	Longitude           *float64
}

