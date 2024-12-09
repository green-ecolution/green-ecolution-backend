package entities

import "time"

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
	ID                 int32                  `json:"id"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	Date               time.Time              `json:"date"`
	Description        string                 `json:"description"`
	Status             WateringStatus         `json:"status"`
	Distance           *float64               `json:"distance"`
	TotalWaterRequired *float64               `json:"total_water_required"`
	Users              []*UserResponse        `json:"users"`
	Treecluster        []*TreeClusterResponse `json:"treecluster"`
	Transporter        *VehicleResponse       `json:"transporter"`
	Trailer            *VehicleResponse       `json:"trailer" validate:"optional"`
} // @Name WateringPlan

type WateringPlanListResponse struct {
	Data       []*WateringPlanResponse `json:"data"`
	Pagination *Pagination             `json:"pagination"`
} // @Name WateringPlanList

type WateringPlanCreateRequest struct {
	Date           time.Time `json:"date"`
	Description    string    `json:"description"`
	TreeclusterIDs []*int32  `json:"tree_cluster_ids"`
	TransporterID  *int32    `json:"transporter_id"`
	TrailerID      *int32    `json:"trailer_id"`
	Users          []*int32  `json:"users_ids"`
} // @Name WateringPlanCreate

type WateringPlanUpdateRequest struct {
	Date           time.Time `json:"date"`
	Description    string    `json:"description"`
	TreeclusterIDs []*int32  `json:"tree_cluster_ids"`
	TransporterID  *int32    `json:"transporter_id"`
	TrailerID      *int32    `json:"trailer_id"`
	Users          []*int32  `json:"users_ids"`
} // @Name WateringPlanUpdate
