package entities

type GeoJSONType string

const (
	FeatureCollection GeoJSONType = "FeatureCollection"
	Feature           GeoJSONType = "Feature"
	LineString        GeoJSONType = "LineString"
)

type GeoJSON struct {
	Type     GeoJSONType      `json:"type"`
	Bbox     []float64        `json:"bbox"`
	Features []GeoJSONFeature `json:"features"`
	Metadata GeoJSONMetadata  `json:"metadata"`
}

type GeoJSONFeature struct {
	Type       GeoJSONType     `json:"type"`
	Bbox       []float64       `json:"bbox"`
	Properties map[string]any  `json:"properties"`
	Geometry   GeoJSONGeometry `json:"geometry"`
}

type GeoJSONGeometry struct {
	Type        GeoJSONType `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

type GeoJSONMetadata struct {
	StartPoint    GeoJSONLocation `json:"start_point"`
	EndPoint      GeoJSONLocation `json:"end_point"`
	WateringPoint GeoJSONLocation `json:"watering_point"`
} // @Name GeoJSONMetadata

type GeoJSONLocation struct {
	Latitude    float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
} // @Name GeoJSONLocation
