package entities

type GeoJsonType string // @Name GeoJsonType

const (
	FeatureCollection GeoJsonType = "FeatureCollection"
)

type GeoJson struct {
	Type     GeoJsonType      `json:"type"`
	Bbox     []float64        `json:"bbox"`
	Features []GeoJsonFeature `json:"features"`
} // @Name GeoJson

type GeoJsonFeature struct {
	Type       GeoJsonType            `json:"type"`
	Bbox       []float64              `json:"bbox"`
	Properties map[string]interface{} `json:"properties"`
	Geometry   GeoJsonGeometry        `json:"geometry"`
} // @Name GeoJsonFeature

type GeoJsonGeometry struct {
	Type        GeoJsonType `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
} // @Name GeoJsonGeometry

type RouteRequest struct {
	VehicleID  int32   `json:"vehicle_id"`
	ClusterIDs []int32 `json:"cluster_ids"`
} // @Name RouteRequest
