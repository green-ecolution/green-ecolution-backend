package mapper

import (
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:TimeToTime
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:TimeToTimePtr
// goverter:extend MapVehicleStatus MapVehicleType
type VehicleHTTPMapper interface {
	FormResponse(*domain.Vehicle) *entities.VehicleResponse
	FromResponseList([]*domain.Vehicle) []*entities.VehicleResponse
}

func MapVehicleStatus(vehicleStatus domain.VehicleStatus) entities.VehicleStatus {
	return entities.VehicleStatus(vehicleStatus)
}

func MapVehicleType(vehicleType domain.VehicleType) entities.VehicleType {
	return entities.VehicleType(vehicleType)
}
