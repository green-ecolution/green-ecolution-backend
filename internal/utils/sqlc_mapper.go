package utils

import (
	"fmt"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func PgTimestampToTime(t pgtype.Timestamp) time.Time {
	return t.Time
}

func PgTimestampToTimePtr(t pgtype.Timestamp) *time.Time {
	if t.Time.IsZero() {
		return nil
	}

	return &t.Time
}

func TimeToPgTimestamp(t *time.Time) pgtype.Timestamp {
	if t == nil {
		return pgtype.Timestamp{}
	}

	return pgtype.Timestamp{
		Time: *t,
	}
}

func PgDateToTime(pgDate pgtype.Date) time.Time {
	if pgDate.Valid {
		return pgDate.Time
	}
	return time.Time{}
}

func TimeToPgDate(date time.Time) (pgtype.Date, error) {
	if date.IsZero() {
		return pgtype.Date{}, fmt.Errorf("invalid date: zero value provided")
	}

	pgDate := pgtype.Date{
		Time:  date,
		Valid: true,
	}

	return pgDate, nil
}

//nolint:gocritic
func ConvertNullableImage(img sqlc.Image) *entities.Image {
	if img.ID == 0 {
		return nil
	}

	return &entities.Image{
		ID:        img.ID,
		CreatedAt: img.CreatedAt.Time,
		UpdatedAt: img.UpdatedAt.Time,
		URL:       img.Url,
	}
}
