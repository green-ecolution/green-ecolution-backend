package mapper

import (
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:PgTimestampToTime
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:PgTimestampToTimePtr
type InternalVehicleRepoMapper interface {
	FromSql(src *sqlc.Vehicle) *entities.Vehicle
	FromSqlList(src []*sqlc.Vehicle) []*entities.Vehicle
}
