package utils

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// P returns a pointer to the value passed as an argument.
func P[T any](v T) *T {
	return &v
}

// RootDir returns the root directory of the project.
func RootDir() string {
	//nolint:dogsled
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b), "../")
	return filepath.Dir(d)
}

// Helper function to decode JSON response
func ParseJSONResponse(body *http.Response, target any) error {
	defer body.Body.Close()
	return json.NewDecoder(body.Body).Decode(target)
}

// Helper function to parse a UUID to a string
func UUIDToString(u uuid.UUID) string {
	if u == uuid.Nil {
		return ""
	}
	return u.String()
}

// Helper function to parse an url to a string
func URLToString(u *url.URL) string {
	if u == nil {
		return ""
	}
	return u.String()
}

func UUIDToPGUUID(userID uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: userID,
		Valid: true,
	}
}

func Scheduler[T any](ctx context.Context, interval time.Duration, process func(ctx context.Context) T) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	if ctx.Err() == nil {
		_ = process(ctx)
	} else {
		slog.Info("Stopping scheduler before first execution due to canceled context")
		return
	}

	for {
		select {
		case <-ticker.C:
			_ = process(ctx)
		case <-ctx.Done():
			slog.Info("Stopping scheduler")
			return
		}
	}
}
