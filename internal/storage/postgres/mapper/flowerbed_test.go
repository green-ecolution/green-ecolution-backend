package mapper_test

import (
	"testing"
	"time"

	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/generated"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

var (
	mapper = &generated.InternalFlowerbedRepoMapperImpl{}
)

func FlowerbedMapper_FromSql(t *testing.T) {
	t.Run("should convert from sql to entity", func(t *testing.T) {
		// given
    timeCreated := time.Now()
    timeUpdated := time.Now()

		src := &sqlc.Flowerbed{
			ID:        1,
			CreatedAt: pgtype.Timestamp{Time: timeCreated},
      UpdatedAt: pgtype.Timestamp{Time: timeUpdated},
		}
		m := &generated.InternalFlowerbedRepoMapperImpl{}

		// when
		got := m.FromSql(src)

		// then
		assert.NotNil(t, got)
	})
}
