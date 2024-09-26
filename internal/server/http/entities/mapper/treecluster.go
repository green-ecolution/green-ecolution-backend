package mapper

import (
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:TimeToTime
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:TimeToTimePtr
// goverter:extend MapWateringStatus MapSoilCondition
type TreeClusterHTTPMapper interface {
  // goverter:ignore Region Trees
	FormResponse(*domain.TreeCluster) *entities.TreeClusterResponse
}

func MapWateringStatus(status domain.TreeClusterWateringStatus) entities.TreeClusterWateringStatus {
	return entities.TreeClusterWateringStatus(status)
}

func MapSoilCondition(condition domain.TreeSoilCondition) entities.TreeSoilCondition {
	return entities.TreeSoilCondition(condition)
}
