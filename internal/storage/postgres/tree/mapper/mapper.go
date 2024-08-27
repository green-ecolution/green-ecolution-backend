package mapper

import (
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:PgTimestampToTime
// goverter:ignoreMissing
type InternalTreeRepoMapper interface {
	FromSqlTree(*sqlc.Tree) *entities.Tree
	FromSqlTreeList([]*sqlc.Tree) []*entities.Tree
}
