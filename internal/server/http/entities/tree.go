package entities

import (
	"time"
)

type TreeResponse struct {
	ID            int32           `json:"id,omitempty"`
	CreatedAt     time.Time       `json:"created_at,omitempty"`
	UpdatedAt     time.Time       `json:"updated_at,omitempty"`
	TreeClusterID *int32          `json:"tree_cluster_id,omitempty" validate:"optional"`
	Sensor        *SensorResponse `json:"sensor,omitempty" validate:"optional"`
	// Images              []*ImageResponse `json:"images,omitempty"`
	Readonly       bool           `json:"readonly,omitempty"`
	PlantingYear   int32          `json:"planting_year,omitempty"`
	Species        string         `json:"species,omitempty"`
	Number         string         `json:"tree_number,omitempty"`
	Latitude       float64        `json:"latitude,omitempty"`
	Longitude      float64        `json:"longitude,omitempty"`
	WateringStatus WateringStatus `json:"watering_status,omitempty"`
	// Description    string         `json:"description,omitempty" validate:"optional"`
} // @Name Tree

type TreeListResponse struct {
	Data       []*TreeResponse `json:"data,omitempty"`
	Pagination Pagination      `json:"pagination,omitempty"`
} // @Name TreeList

type TreeCreateRequest struct {
	TreeClusterID *int32  `json:"tree_cluster_id,omitempty" validate:"optional"`
	Readonly      bool    `json:"readonly,omitempty"`
	PlantingYear  int32   `json:"planting_year,omitempty"`
	Species       string  `json:"species,omitempty"`
	Number        string  `json:"tree_number,omitempty"`
	Latitude      float64 `json:"latitude,omitempty"`
	Longitude     float64 `json:"longitude,omitempty"`
	SensorID      *int32  `json:"sensor_id,omitempty" validate:"optional"`
	Description   string  `json:"description,omitempty" validate:"optional"`
} // @Name TreeCreate

type TreeUpdateRequest struct {
	TreeClusterID *int32  `json:"tree_cluster_id,omitempty" validate:"optional"`
	Readonly      bool    `json:"readonly,omitempty"`
	PlantingYear  int32   `json:"planting_year,omitempty"`
	Species       string  `json:"species,omitempty"`
	Number        string  `json:"tree_number,omitempty"`
	Latitude      float64 `json:"latitude,omitempty"`
	Longitude     float64 `json:"longitude,omitempty"`
	SensorID      *int32  `json:"sensor_id,omitempty" validate:"optional"`
	Description   string  `json:"description,omitempty" validate:"optional"`
} // @Name TreeUpdate

type TreeAddImagesRequest struct {
	ImageIDs []*int32 `json:"image_ids,omitempty"`
} // @Name TreeAddImages

type TreeAddSensorRequest struct {
	SensorID *int32 `json:"sensor_id,omitempty"`
} // @Name TreeAddSensor
