package mapper

import (
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:TimeToTime
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:TimeToTimePtr
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:UUIDToString
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:URLToString
// goverter:extend MapWateringPlanStatus MapVehicleStatus MapVehicleType MapDrivingLicense
// goverter:extend MapWateringStatus MapSensorStatus MapSoilCondition MapTreesToIDs
type WateringPlanHTTPMapper interface {
	FromResponse(*domain.WateringPlan) *entities.WateringPlanResponse
	FromResponseList([]*domain.WateringPlan) []*entities.WateringPlanResponse
	FromCreateRequest(*entities.WateringPlanCreateRequest) *domain.WateringPlanCreate
	FromUpdateRequest(*entities.WateringPlanUpdateRequest) *domain.WateringPlanUpdate

	// goverter:map Trees TreeIDs
	FromInListResponse(*domain.TreeCluster) *entities.TreeClusterInListResponse
}

func MapWateringPlanStatus(status domain.WateringPlanStatus) entities.WateringPlanStatus {
	return entities.WateringPlanStatus(status)
}
