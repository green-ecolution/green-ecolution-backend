package entities

import "time"

type TreeSoilCondition string

const (
	TreeSoilConditionSchluffig TreeSoilCondition = "schluffig"
	TreeSoilConditionSandig    TreeSoilCondition = "sandig"
	TreeSoilConditionLehmig    TreeSoilCondition = "lehmig"
	TreeSoilConditionTonig     TreeSoilCondition = "tonig"
	TreeSoilConditionUnknown   TreeSoilCondition = "unknown"
)

type TreeCluster struct {
	ID             int32
	CreatedAt      time.Time
	UpdatedAt      time.Time
	WateringStatus WateringStatus
	LastWatered    *time.Time
	MoistureLevel  float64
	Region         *Region
	Address        string
	Description    string
	Archived       bool
	Latitude       *float64
	Longitude      *float64
	Trees          []*Tree
	SoilCondition  TreeSoilCondition
	Name           string
}

type TreeClusterCreate struct {
	Address       string
	Description   string
	Name          string `validate:"required"`
	SoilCondition TreeSoilCondition
	TreeIDs       []*int32
}

type TreeClusterUpdate struct {
	Address       string
	Description   string
	SoilCondition TreeSoilCondition
	TreeIDs       []*int32
	Name          string `validate:"required"`
}
