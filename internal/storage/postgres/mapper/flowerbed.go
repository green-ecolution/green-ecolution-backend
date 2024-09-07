package mapper

import (
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:PgTimestampToTime
type InternalFlowerbedRepoMapper interface {
	// goverter:ignore Sensor Images
	FromSql(src *sqlc.Flowerbed) *entities.Flowerbed
	FromSqlList(src []*sqlc.Flowerbed) []*entities.Flowerbed
}
