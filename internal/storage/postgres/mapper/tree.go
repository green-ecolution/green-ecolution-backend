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
type InternalTreeRepoMapper interface {
	// goverter:ignore Sensor Images TreeCluster
	// goverter:map AdditionalInformations AdditionalInfo | github.com/green-ecolution/green-ecolution-backend/internal/utils:MapAdditionalInfo
	FromSql(*sqlc.Tree) (*entities.Tree, error)
	FromSqlList([]*sqlc.Tree) ([]*entities.Tree, error)
}
