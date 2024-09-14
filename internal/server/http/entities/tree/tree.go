package tree

import (
	"time"
)

type TreeResponse struct {
	ID            int32                  `json:"id,omitempty"`
	CreatedAt     time.Time              `json:"created_at,omitempty"`
	UpdatedAt     time.Time              `json:"updated_at,omitempty"`
	TreeClusterID *int32                 `json:"tree_cluster_id,omitempty"`
	// Sensor        *sensor.SensorResponse `json:"sensor,omitempty"`
	// Images              []*ImageResponse `json:"images,omitempty"`
	Age                 int32   `json:"age,omitempty"`
	HeightAboveSeaLevel float64 `json:"height_above_sea_level,omitempty"`
	PlantingYear        int32   `json:"planting_year,omitempty"`
	Species             string  `json:"species,omitempty"`
	Number              int32   `json:"number,omitempty"`
	Latitude            float64 `json:"latitude,omitempty"`
	Longitude           float64 `json:"longitude,omitempty"`
} // @Name Tree
