package mapper

import (
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:PgTimestampToTime
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:PgTimestampToTimePtr
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:PgDateToTime
// goverter:extend MapWateringPlanStatus
type InternalWateringPlanRepoMapper interface {
	// goverter:ignore Users TreeClusters Transporter Trailer ConsumedWaterList
	FromSql(src *sqlc.WateringPlan) *entities.WateringPlan
	FromSqlList(src []*sqlc.WateringPlan) []*entities.WateringPlan
}

func MapWateringPlanStatus(wateringPlanStatus sqlc.WateringPlanStatus) entities.WateringPlanStatus {
	return entities.WateringPlanStatus(wateringPlanStatus)
}
