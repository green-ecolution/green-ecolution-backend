package entities

type GeoJSONType string

const (
	FeatureCollection GeoJSONType = "FeatureCollection"
)

type GeoJSON struct {
	Type     GeoJSONType
	Bbox     []float64
	Features []map[string]interface{}
}
