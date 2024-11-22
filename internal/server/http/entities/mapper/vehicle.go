package mapper

import (
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:TimeToTime
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:TimeToTimePtr
// goverter:extend MapVehicleStatus MapVehicleType MapVehicleStatusReq MapVehicleTypeReq MapDriverLicence MapDriverLicenceReq
type VehicleHTTPMapper interface {
	FormResponse(*domain.Vehicle) *entities.VehicleResponse
	FromResponseList([]*domain.Vehicle) []*entities.VehicleResponse
	FromCreateRequest(*entities.VehicleCreateRequest) *domain.VehicleCreate
	FromUpdateRequest(*entities.VehicleUpdateRequest) *domain.VehicleUpdate
}

func MapVehicleStatus(vehicleStatus domain.VehicleStatus) entities.VehicleStatus {
	return entities.VehicleStatus(vehicleStatus)
}

func MapVehicleType(vehicleType domain.VehicleType) entities.VehicleType {
	return entities.VehicleType(vehicleType)
}

func MapDriverLicence(driverLicence domain.DriverLicence) entities.DriverLicence {
	return entities.DriverLicence(driverLicence)
}

func MapVehicleStatusReq(vehicleStatus entities.VehicleStatus) domain.VehicleStatus {
	return domain.VehicleStatus(vehicleStatus)
}

func MapVehicleTypeReq(vehicleType entities.VehicleType) domain.VehicleType {
	return domain.VehicleType(vehicleType)
}

func MapDriverLicenceReq(driverLicence entities.DriverLicence) domain.DriverLicence {
	return domain.DriverLicence(driverLicence)
}
