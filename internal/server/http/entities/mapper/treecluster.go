package mapper

import (
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:TimeToTime
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:TimeToTimePtr
// goverter:extend MapWateringStatus MapSoilCondition MapSoilConditionReq
type TreeClusterHTTPMapper interface {
	// goverter:ignore Region Trees
	FromResponse(*domain.TreeCluster) *entities.TreeClusterResponse

	FromCreateRequest(*entities.TreeClusterCreateRequest) *domain.TreeClusterCreate
	FromUpdateRequest(*entities.TreeClusterUpdateRequest) *domain.TreeClusterUpdate
}

func MapWateringStatus(status domain.WateringStatus) entities.WateringStatus {
	return entities.WateringStatus(status)
}

func MapSoilCondition(condition domain.TreeSoilCondition) entities.TreeSoilCondition {
	return entities.TreeSoilCondition(condition)
}

func MapSoilConditionReq(condition entities.TreeSoilCondition) domain.TreeSoilCondition {
	return domain.TreeSoilCondition(condition)
}
