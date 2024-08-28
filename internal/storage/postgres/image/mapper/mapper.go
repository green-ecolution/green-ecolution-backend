package mapper

import (
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:PgTimestampToTime
// goverter:ignoreMissing
type InternalImageRepoMapper interface {
	// goverter:map Url URL
	FromSql(src *sqlc.Image) *entities.Image
	FromSqlList(src []*sqlc.Image) []*entities.Image
}
