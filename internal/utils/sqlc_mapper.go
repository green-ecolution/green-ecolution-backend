package utils

import (
	"encoding/json"
	"fmt"
	"time"

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
		Time:  *t,
		Valid: true,
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

func MapAdditionalInfo(src []byte) (map[string]any, error) {
	if len(src) == 0 {
		return nil, nil
	}

	additionalInfo := make(map[string]any, 0)
	err := json.Unmarshal(src, &additionalInfo)
	if err != nil {
		return nil, err
	}
	return additionalInfo, nil
}

func MapAdditionalInfoToByte(src map[string]any) ([]byte, error) {
	if src == nil {
		return nil, nil
	}

	if len(src) == 0 {
		return nil, nil
	}

	additionalInfo, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}

	return additionalInfo, nil
}
