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

func TestPgDateToTime(t *testing.T) {
	t.Run("should return time from pgtype.Date", func(t *testing.T) {
		date := pgtype.Date{Time: time.Now(), Valid: true}
		result := PgDateToTime(date)

		assert.Equal(t, date.Time, result)
	})

	t.Run("should return zero time for zero pgtype.Date", func(t *testing.T) {
		date := pgtype.Date{Valid: false}
		result := PgDateToTime(date)

		assert.True(t, result.IsZero())
	})
}

func TestTimeToPgDate(t *testing.T) {
	t.Run("Convert current time", func(t *testing.T) {
		date := time.Now()

		pgDate, err := TimeToPgDate(date)
		assert.NoError(t, err)

		assert.Equal(t, date.Year(), pgDate.Time.Year())
		assert.Equal(t, date.Month(), pgDate.Time.Month())
		assert.Equal(t, date.Day(), pgDate.Time.Day())
		assert.Equal(t, date.Hour(), pgDate.Time.Hour())
		assert.Equal(t, date.Minute(), pgDate.Time.Minute())
		assert.Equal(t, date.Second(), pgDate.Time.Second())
		assert.Equal(t, date.Nanosecond(), pgDate.Time.Nanosecond())
	})

	t.Run("Convert empty time", func(t *testing.T) {
		date := time.Time{}

		pgDate, err := TimeToPgDate(date)

		assert.Error(t, err)
		assert.Empty(t, pgDate)
		assert.Equal(t, "invalid date: zero value provided", err.Error())
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
