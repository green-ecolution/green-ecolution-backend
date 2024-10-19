package entities

import (
	"time"
)

type TreeSoilCondition string // @Name SoilCondition

const (
	TreeSoilConditionSchluffig TreeSoilCondition = "schluffig"
	TreeSoilConditionSandig    TreeSoilCondition = "sandig"
	TreeSoilConditionLehmig    TreeSoilCondition = "lehmig"
	TreeSoilConditionTonig     TreeSoilCondition = "tonig"
	TreeSoilConditionUnknown   TreeSoilCondition = "unknown"
)

type TreeClusterResponse struct {
	ID             int32             `json:"id"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	WateringStatus WateringStatus    `json:"watering_status"`
	LastWatered    *time.Time        `json:"last_watered,omitempty" validate:"optional"`
	MoistureLevel  float64           `json:"moisture_level"`
	Region         *RegionResponse   `json:"region,omitempty" validate:"optional"`
	Address        string            `json:"address"`
	Description    string            `json:"description"`
	Archived       bool              `json:"archived"`
	Latitude       *float64          `json:"latitude"`
	Longitude      *float64          `json:"longitude"`
	Trees          []*TreeResponse   `json:"trees" validate:"optional"`
	SoilCondition  TreeSoilCondition `json:"soil_condition"`
	Name           string            `json:"name"`
} // @Name TreeCluster

type TreeClusterListResponse struct {
	Data       []*TreeClusterResponse `json:"data"`
	Pagination *Pagination            `json:"pagination"`
} // @Name TreeClusterList

type TreeClusterCreateRequest struct {
	Address       string            `json:"address"`
	Description   string            `json:"description"`
	TreeIDs       []*int32          `json:"tree_ids"`
	SoilCondition TreeSoilCondition `json:"soil_condition"`
	Name          string            `json:"name"`
} // @Name TreeClusterCreate

type TreeClusterUpdateRequest struct {
	Address       string            `json:"address"`
	Description   string            `json:"description"`
	TreeIDs       []*int32          `json:"tree_ids"`
	SoilCondition TreeSoilCondition `json:"soil_condition"`
	Name          string            `json:"name"`
} // @Name TreeClusterUpdate

type TreeClusterAddTreesRequest struct {
	TreeIDs []*int32 `json:"tree_ids"`
} // @Name TreeClusterAddTrees
