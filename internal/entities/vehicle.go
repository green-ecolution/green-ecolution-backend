package entities

import (
	"time"
)

type VehicleStatus string

const (
	VehicleStatusActive       VehicleStatus = "active"
	VehicleStatusAvailable    VehicleStatus = "available"
	VehicleStatusNotAvailable VehicleStatus = "not available"
	VehicleStatusUnknown      VehicleStatus = "unknown"
)

type VehicleType string

const (
	VehicleTypeTransporter VehicleType = "transporter"
	VehicleTypeTrailer     VehicleType = "trailer"
	VehicleTypeUnknown     VehicleType = "unknown"
)

type Vehicle struct {
	ID            int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
	NumberPlate   string
	Description   string
	WaterCapacity float64
	Status        VehicleStatus
	Type          VehicleType
}

type VehicleCreate struct {
	NumberPlate   string `validate:"required"`
	Description   string
	WaterCapacity float64 `validate:"gt=0"`
	Status        VehicleStatus
	Type          VehicleType
}

type VehicleUpdate struct {
	NumberPlate   string `validate:"required"`
	Description   string 
	WaterCapacity float64 `validate:"gt=0"`
	Status        VehicleStatus
	Type          VehicleType
}
