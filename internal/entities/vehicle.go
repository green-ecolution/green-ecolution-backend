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
	Model         string
	DriverLicence DriverLicence
	Height        float64
	Width         float64
	Length        float64
}

type VehicleCreate struct {
	NumberPlate   string `validate:"required"`
	Description   string
	WaterCapacity float64       `validate:"gt=0"`
	Status        VehicleStatus `validate:"oneof=active available 'not available' unknown"`
	Type          VehicleType   `validate:"oneof=transporter trailer unknown"`
	Model         string        `validate:"required"`
	DriverLicence DriverLicence `validate:"oneof=B BE C"`
	Height        float64       `validate:"gt=0"`
	Width         float64       `validate:"gt=0"`
	Length        float64       `validate:"gt=0"`
}

type VehicleUpdate struct {
	NumberPlate   string `validate:"required"`
	Description   string
	WaterCapacity float64       `validate:"gt=0"`
	Status        VehicleStatus `validate:"oneof=active available 'not available' unknown"`
	Type          VehicleType   `validate:"oneof=transporter trailer unknown"`
	Model         string        `validate:"required"`
	DriverLicence DriverLicence `validate:"oneof=B BE C"`
	Height        float64       `validate:"gt=0"`
	Width         float64       `validate:"gt=0"`
	Length        float64       `validate:"gt=0"`
}
