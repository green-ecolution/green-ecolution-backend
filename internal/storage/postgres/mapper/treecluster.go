package mapper

import (
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:PgTimestampToTime
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:PgTimestampToTimePtr
// goverter:extend MapWateringStatus MapSoilCondition
// goverter:ignoreMissing
type InternalTreeClusterRepoMapper interface {
	FromSql(*sqlc.TreeCluster) *entities.TreeCluster
	FromSqlList([]*sqlc.TreeCluster) []*entities.TreeCluster
}

func MapWateringStatus(status sqlc.TreeClusterWateringStatus) entities.TreeClusterWateringStatus {
	return entities.TreeClusterWateringStatus(status)
}

func MapSoilCondition(condition sqlc.TreeSoilCondition) entities.TreeSoilCondition {
	return entities.TreeSoilCondition(condition)
}
