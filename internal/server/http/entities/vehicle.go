package entities

import "time"

type VehicleStatus string // @Name VehicleStatus

const (
	VehicleStatusActive       VehicleStatus = "active"
	VehicleStatusAvailable    VehicleStatus = "available"
	VehicleStatusNotAvailable VehicleStatus = "not available"
	VehicleStatusUnknown      VehicleStatus = "unknown"
)

type VehicleType string // @Name VehicleType

const (
	VehicleTypeTransporter VehicleType = "transporter"
	VehicleTypeTrailer     VehicleType = "trailer"
	VehicleTypeUnknown     VehicleType = "unknown"
)

type VehicleResponse struct {
	ID                   int32         `json:"id"`
	CreatedAt            time.Time     `json:"created_at"`
	UpdatedAt            time.Time     `json:"updated_at"`
	NumberPlate          string        `json:"number_plate"`
	Description          string        `json:"description"`
	WaterCapacity        float64       `json:"water_capacity"`
	Status               VehicleStatus `json:"status"`
	Type                 VehicleType   `json:"type"`
	Model                string        `json:"model"`
	DriverLicence DriverLicence `json:"driver_licence"`
	Height               float64       `json:"height"`
	Width                float64       `json:"width"`
	Length               float64       `json:"length"`
} // @Name Vehicle

type VehicleListResponse struct {
	Data       []*VehicleResponse `json:"data"`
	Pagination *Pagination        `json:"pagination"`
} // @Name VehicleList

type VehicleCreateRequest struct {
	NumberPlate          string        `json:"number_plate"`
	Description          string        `json:"description"`
	WaterCapacity        float64       `json:"water_capacity"`
	Status               VehicleStatus `json:"status"`
	Type                 VehicleType   `json:"type"`
	Model                string        `json:"model"`
	DriverLicence DriverLicence `json:"driver_licence"`
	Height               float64       `json:"height"`
	Width                float64       `json:"width"`
	Length               float64       `json:"length"`
} // @Name VehicleCreate

type VehicleUpdateRequest struct {
	NumberPlate          string        `json:"number_plate"`
	Description          string        `json:"description"`
	WaterCapacity        float64       `json:"water_capacity"`
	Status               VehicleStatus `json:"status"`
	Type                 VehicleType   `json:"type"`
	Model                string        `json:"model"`
	DriverLicence DriverLicence `json:"driver_licence"`
	Height               float64       `json:"height"`
	Width                float64       `json:"width"`
	Length               float64       `json:"length"`
} // @Name VehicleUpdate
