package mapper

import (
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:PgTimestampToTime
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:PgTimestampToTimePtr
// goverter:extend MapWateringStatus StringPtrToString
type InternalTreeRepoMapper interface {
	// goverter:ignore Sensor Images TreeCluster
	// goverter:map TreeNumber Number
	FromSql(*sqlc.Tree) *entities.Tree
	FromSqlList([]*sqlc.Tree) []*entities.Tree
}

func StringPtrToString(source *string) string {
	if source == nil {
		return ""
	}
	return *source
}
