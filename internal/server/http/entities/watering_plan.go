package entities

import (
	"time"

	"github.com/google/uuid"
)

type WateringPlanStatus string // @Name WateringPlanStatus

const (
	WateringPlanStatusPlanned     WateringPlanStatus = "planned"
	WateringPlanStatusActive      WateringPlanStatus = "active"
	WateringPlanStatusCanceled    WateringPlanStatus = "canceled"
	WateringPlanStatusFinished    WateringPlanStatus = "finished"
	WateringPlanStatusNotCompeted WateringPlanStatus = "not competed"
	WateringPlanStatusUnknown     WateringPlanStatus = "unknown"
)

type WateringPlanResponse struct {
	ID                 int32                        `json:"id"`
	CreatedAt          time.Time                    `json:"created_at"`
	UpdatedAt          time.Time                    `json:"updated_at"`
	Date               time.Time                    `json:"date"`
	Description        string                       `json:"description"`
	Status             WateringPlanStatus           `json:"status"`
	Distance           *float64                     `json:"distance"`
	TotalWaterRequired *float64                     `json:"total_water_required"`
	UserIDs            []*uuid.UUID                 `json:"user_ids"`
	TreeClusters       []*TreeClusterInListResponse `json:"treeclusters"`
	Transporter        *VehicleResponse             `json:"transporter"`
	Trailer            *VehicleResponse             `json:"trailer" validate:"optional"`
	CancellationNote   string                       `json:"cancellation_note"`
	Evaluation         []*EvaluationValue           `json:"evaluation"`
} // @Name WateringPlan

type WateringPlanInListResponse struct {
	ID                 int32                        `json:"id"`
	CreatedAt          time.Time                    `json:"created_at"`
	UpdatedAt          time.Time                    `json:"updated_at"`
	Date               time.Time                    `json:"date"`
	Description        string                       `json:"description"`
	Status             WateringPlanStatus           `json:"status"`
	Distance           *float64                     `json:"distance"`
	TotalWaterRequired *float64                     `json:"total_water_required"`
	UserIDs            []*uuid.UUID                 `json:"user_ids"`
	TreeClusters       []*TreeClusterInListResponse `json:"treeclusters"`
	Transporter        *VehicleResponse             `json:"transporter"`
	Trailer            *VehicleResponse             `json:"trailer" validate:"optional"`
	CancellationNote   string                       `json:"cancellation_note"`
} // @Name WateringPlanInList

type WateringPlanListResponse struct {
	Data       []*WateringPlanInListResponse `json:"data"`
	Pagination *Pagination                   `json:"pagination"`
} // @Name WateringPlanList

type WateringPlanCreateRequest struct {
	Date           time.Time `json:"date"`
	Description    string    `json:"description"`
	TreeClusterIDs []*int32  `json:"tree_cluster_ids"`
	TransporterID  *int32    `json:"transporter_id"`
	TrailerID      *int32    `json:"trailer_id" validate:"optional"`
	UserIDs        []string  `json:"user_ids"`
} // @Name WateringPlanCreate

type WateringPlanUpdateRequest struct {
	Date             time.Time          `json:"date"`
	Description      string             `json:"description"`
	TreeClusterIDs   []*int32           `json:"tree_cluster_ids"`
	TransporterID    *int32             `json:"transporter_id"`
	TrailerID        *int32             `json:"trailer_id" validate:"optional"`
	UserIDs          []string           `json:"user_ids"`
	CancellationNote string             `json:"cancellation_note"`
	Status           WateringPlanStatus `json:"status"`
	Evaluation       []*EvaluationValue `json:"evaluation" validate:"optional"`
} // @Name WateringPlanUpdate

type EvaluationValue struct {
	WateringPlanID int32    `json:"watering_plan_id"`
	TreeClusterID  int32    `json:"tree_cluster_id"`
	ConsumedWater  *float64 `json:"consumed_water"`
} // @Name EvaluationValue
