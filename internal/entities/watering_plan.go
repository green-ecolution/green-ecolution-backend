package entities

import (
	"time"
)

type WateringPlanStatus string

const (
	WateringPlanStatusPlanned     WateringPlanStatus = "planned"
	WateringPlanStatusActive      WateringPlanStatus = "active"
	WateringPlanStatusCancelled   WateringPlanStatus = "cancelled"
	WateringPlanStatusFinished    WateringPlanStatus = "finished"
	WateringPlanStatusNotCompeted WateringPlanStatus = "not competed"
	WateringPlanStatusUnknown     WateringPlanStatus = "unknown"
)

type WateringPlan struct {
	ID        int32
	CreatedAt time.Time
	UpdatedAt time.Time
	Date         time.Time
	Description        string
	WateringPlanStatus WateringPlanStatus
	Distance           *float64
	TotalWaterRequired *float64
	Users              []*User
	Treecluster        []*TreeCluster
	Transporter			*Vehicle
	Trailer				*Vehicle 
}
