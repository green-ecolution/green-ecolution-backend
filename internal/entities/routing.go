package entities

import (
	"time"
)

type Routing struct {
	EstimatedTime   time.Time
	EstimatedLength int
}

type GeoJSON struct {
	Type     *string
	Bbox     []float64
	Features []map[string]interface{}
}
