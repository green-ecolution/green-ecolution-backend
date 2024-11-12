package utils

import (
	"testing"
	"time"

	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestPgTimestampToTime(t *testing.T) {
	t.Run("should return time from pgtype.Timestamp", func(t *testing.T) {
		ts := pgtype.Timestamp{Time: time.Now()}
		result := PgTimestampToTime(ts)
		assert.Equal(t, ts.Time, result)
	})

	t.Run("should return zero time for zero pgtype.Timestamp", func(t *testing.T) {
		ts := pgtype.Timestamp{}
		result := PgTimestampToTime(ts)
		assert.True(t, result.IsZero())
	})
}

func TestPgTimestampToTimePtr(t *testing.T) {
	t.Run("should return time pointer from pgtype.Timestamp", func(t *testing.T) {
		ts := pgtype.Timestamp{Time: time.Now()}
		result := PgTimestampToTimePtr(ts)
		assert.NotNil(t, result)
		assert.Equal(t, ts.Time, *result)
	})

	t.Run("should return nil for zero pgtype.Timestamp", func(t *testing.T) {
		ts := pgtype.Timestamp{}
		result := PgTimestampToTimePtr(ts)
		assert.Nil(t, result)
	})

	t.Run("should return pointer to non-zero time", func(t *testing.T) {
		now := time.Now()
		ts := pgtype.Timestamp{Time: now}
		result := PgTimestampToTimePtr(ts)
		assert.NotNil(t, result)
		assert.Equal(t, now, *result)
	})
}

func TestTimeToPgTimestamp(t *testing.T) {
	t.Run("should return pgtype.Timestamp from time pointer", func(t *testing.T) {
		tm := time.Now()
		result := TimeToPgTimestamp(&tm)
		assert.Equal(t, tm, result.Time)
	})

	t.Run("should return zero pgtype.Timestamp for nil time pointer", func(t *testing.T) {
		result := TimeToPgTimestamp(nil)
		assert.True(t, result.Time.IsZero())
	})

	t.Run("should convert current time to pgtype.Timestamp", func(t *testing.T) {
		now := time.Now()
		result := TimeToPgTimestamp(&now)
		assert.Equal(t, now, result.Time)
	})
}

func TestConvertNullableImage(t *testing.T) {
	t.Run("should convert sqlc.Image to entities.Image", func(t *testing.T) {
		img := sqlc.Image{
			ID:        1,
			CreatedAt: pgtype.Timestamp{Time: time.Now()},
			UpdatedAt: pgtype.Timestamp{Time: time.Now()},
			Url:       "http://example.com/image.jpg",
		}
		result := ConvertNullableImage(img)
		assert.NotNil(t, result)
		assert.Equal(t, img.ID, result.ID)
		assert.Equal(t, img.CreatedAt.Time, result.CreatedAt)
		assert.Equal(t, img.UpdatedAt.Time, result.UpdatedAt)
		assert.Equal(t, img.Url, result.URL)
	})

	t.Run("should return nil for sqlc.Image with zero ID", func(t *testing.T) {
		img := sqlc.Image{
			ID:        0,
			CreatedAt: pgtype.Timestamp{Time: time.Now()},
			UpdatedAt: pgtype.Timestamp{Time: time.Now()},
			Url:       "http://example.com/image.jpg",
		}
		result := ConvertNullableImage(img)
		assert.Nil(t, result)
	})

	t.Run("should handle sqlc.Image with zero CreatedAt and UpdatedAt", func(t *testing.T) {
		img := sqlc.Image{
			ID:        1,
			CreatedAt: pgtype.Timestamp{},
			UpdatedAt: pgtype.Timestamp{},
			Url:       "http://example.com/image.jpg",
		}
		result := ConvertNullableImage(img)
		assert.NotNil(t, result)
		assert.Equal(t, img.ID, result.ID)
		assert.True(t, result.CreatedAt.IsZero())
		assert.True(t, result.UpdatedAt.IsZero())
		assert.Equal(t, img.Url, result.URL)
	})
}

