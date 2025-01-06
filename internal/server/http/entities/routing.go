package entities

type GeoJSONType string // @Name GeoJsonType

const (
	FeatureCollection GeoJSONType = "FeatureCollection"
)

type GeoJSON struct {
	Type     GeoJSONType              `json:"type,omitempty"`
	Bbox     []float64                `json:"bbox,omitempty"`
	Features []map[string]interface{} `json:"features,omitempty"`
} // @Name GeoJson

type RouteRequest struct {
	VehicleID  int32   `json:"vehicle_id"`
	ClusterIDs []int32 `json:"cluster_ids"`
} // @Name RouteRequest
