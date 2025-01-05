package entities

type GeoJSON struct {
	Type     *string                  `json:"type,omitempty"`
	Bbox     []float64                `json:"bbox,omitempty"`
	Features []map[string]interface{} `json:"features,omitempty"`
} // @Name GeoJson

type RouteRequest struct {
	VehicleID  int32   `json:"vehicle_id"`
	ClusterIDs []int32 `json:"cluster_ids"`
} // @Name RouteRequest
