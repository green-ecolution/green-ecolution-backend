package entities

type GeoJsonType string

const (
	FeatureCollection GeoJsonType = "FeatureCollection"
	Feature           GeoJsonType = "Feature"
	LineString        GeoJsonType = "LineString"
)

type GeoJson struct {
	Type     GeoJsonType      `json:"type"`
	Bbox     []float64        `json:"bbox"`
	Features []GeoJsonFeature `json:"features"`
}

type GeoJsonFeature struct {
	Type       GeoJsonType            `json:"type"`
	Bbox       []float64              `json:"bbox"`
	Properties map[string]interface{} `json:"properties"`
	Geometry   GeoJsonGeometry        `json:"geometry"`
}

type GeoJsonGeometry struct {
	Type        GeoJsonType `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}
