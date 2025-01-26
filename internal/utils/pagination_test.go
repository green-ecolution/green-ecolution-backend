package utils

import (
	"context"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/stretchr/testify/assert"
)

func TestCreatePagination(t *testing.T) {
	t.Run("should return a valid pagination", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "page", int32(2))
		ctx = context.WithValue(ctx, "limit", int32(10))

		totalCount := int64(50)
		pagination := CreatePagination(ctx, totalCount)

		expectedNextPage := int32(3)
		expectedPrevPage := int32(1)
		expectedPagination := &entities.Pagination{
			Total:       totalCount,
			CurrentPage: 2,
			TotalPages:  5,
			NextPage:    &expectedNextPage,
			PrevPage:    &expectedPrevPage,
		}

		// assert pagination values
		assert.NotNil(t, pagination)
		assert.Equal(t, expectedPagination.Total, pagination.Total)
		assert.Equal(t, expectedPagination.CurrentPage, pagination.CurrentPage)
		assert.Equal(t, expectedPagination.TotalPages, pagination.TotalPages)
		assert.Equal(t, expectedPagination.NextPage, pagination.NextPage)
		assert.Equal(t, expectedPagination.PrevPage, pagination.PrevPage)
	})

	t.Run("should return nil on limit with value -1", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "page", int32(1))
		ctx = context.WithValue(ctx, "limit", int32(-1))

		pagination := CreatePagination(ctx, int64(50))

		assert.Nil(t, pagination)
	})

	t.Run("should return nil on no context values", func(t *testing.T) {
		pagination := CreatePagination(context.Background(), int64(50))

		assert.Nil(t, pagination)
	})
}
