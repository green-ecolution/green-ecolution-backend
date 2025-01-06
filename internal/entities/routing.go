package entities

import (
	"time"
)

type GeoJSONType string

const (
	FeatureCollection GeoJSONType = "FeatureCollection"
)

type Routing struct {
	EstimatedTime   time.Time
	EstimatedLength int
}

type GeoJSON struct {
	Type     GeoJSONType
	Bbox     []float64
	Features []map[string]interface{}
}
