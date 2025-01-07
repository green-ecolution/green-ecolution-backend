package entities

type GeoJSONType string // @Name GeoJsonType

const (
	FeatureCollection GeoJSONType = "FeatureCollection"
)

type GeoJSON struct {
	Type     GeoJSONType      `json:"type"`
	Bbox     []float64        `json:"bbox"`
	Features []GeoJSONFeature `json:"features"`
} // @Name GeoJson

type GeoJSONFeature struct {
	Type       GeoJSONType     `json:"type"`
	Bbox       []float64       `json:"bbox"`
	Properties map[string]any  `json:"properties"`
	Geometry   GeoJSONGeometry `json:"geometry"`
} // @Name GeoJsonFeature

type GeoJSONGeometry struct {
	Type        GeoJSONType `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
} // @Name GeoJsonGeometry

type RouteRequest struct {
	TransporterID int32   `json:"transporter_id"`
	TrailerID     *int32  `json:"trailer_id,omitempty" validate:"optional"`
	ClusterIDs    []int32 `json:"cluster_ids"`
} // @Name RouteRequest
