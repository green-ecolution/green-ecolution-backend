package entities

import "time"

type TreeClusterWateringStatus string

const (
	TreeClusterWateringStatusGood     TreeClusterWateringStatus = "good"
	TreeClusterWateringStatusModerate TreeClusterWateringStatus = "moderate"
	TreeClusterWateringStatusBad      TreeClusterWateringStatus = "bad"
	TreeClusterWateringStatusUnknown  TreeClusterWateringStatus = "unknown"
)

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
	WateringStatus TreeClusterWateringStatus
	LastWatered    *time.Time
	MoistureLevel  float64
	Region         *Region
	Address        string
	Description    string
	Archived       bool
	Latitude       float64
	Longitude      float64
	Trees          []*Tree
	SoilCondition  TreeSoilCondition
	Name           string
}
