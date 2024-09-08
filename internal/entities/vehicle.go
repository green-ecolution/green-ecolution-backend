package entities

import (
	"time"
)

type Vehicle struct {
	ID            int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
	NumberPlate   string
	Description   string
	WaterCapacity float64
}

type CreateVehicle struct {
	NumberPlate   string
	Description   string
	WaterCapacity float64
}

type UpdateVehicle struct {
	ID            int32
	NumberPlate   *string
	Description   *string
	WaterCapacity *float64
}
