package mapper

import (
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:PgTimestampToTime
type InternalDepartureMapper interface {
	FromSql(src *sqlc.Departure) *entities.Departure
	FromSqlList(src []*sqlc.Departure) []*entities.Departure
}
