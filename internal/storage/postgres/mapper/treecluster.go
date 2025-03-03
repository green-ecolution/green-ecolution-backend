package mapper

import (
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:PgTimestampToTime
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:PgTimestampToTimePtr
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:StringPtrToString
// goverter:extend MapWateringStatus MapSoilCondition
// goverter:ignoreMissing
type InternalTreeClusterRepoMapper interface {
	// goverter:map AdditionalInformations AdditionalInfo | github.com/green-ecolution/green-ecolution-backend/internal/utils:MapAdditionalInfo
	FromSql(*sqlc.TreeCluster) (*entities.TreeCluster, error)
	FromSqlList([]*sqlc.TreeCluster) ([]*entities.TreeCluster, error)

	FromSqlRegionWithCount(src *sqlc.GetAllTreeClusterRegionsWithWateringPlanCountRow) (*entities.RegionEvaluation, error)
	FromSqlRegionListWithCount(src []*sqlc.GetAllTreeClusterRegionsWithWateringPlanCountRow) ([]*entities.RegionEvaluation, error)
}

func FromSqlRegionWithCount(src *sqlc.GetAllTreeClusterRegionsWithWateringPlanCountRow) (*entities.RegionEvaluation, error) {
	if src == nil {
		return nil, nil
	}

	return &entities.RegionEvaluation{
		Name:              src.RegionName,
		WateringPlanCount: int64(src.WateringPlanCount),
	}, nil
}

func FromSqlRegionListWithCount(src []*sqlc.GetAllTreeClusterRegionsWithWateringPlanCountRow) ([]*entities.RegionEvaluation, error) {
	var result []*entities.RegionEvaluation
	for _, v := range src {
		mapped, err := FromSqlRegionWithCount(v)
		if err != nil {
			return nil, err
		}
		result = append(result, mapped)
	}
	return result, nil
}

func MapWateringStatus(status sqlc.WateringStatus) entities.WateringStatus {
	return entities.WateringStatus(status)
}

func MapSoilCondition(condition sqlc.TreeSoilCondition) entities.TreeSoilCondition {
	return entities.TreeSoilCondition(condition)
}
