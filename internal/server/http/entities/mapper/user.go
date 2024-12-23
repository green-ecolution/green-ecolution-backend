package mapper

import (
	"net/url"

	"github.com/google/uuid"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:TimeToTime
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:UUIDToString
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:URLToString
// goverter:extend UUIDToUUID
// goverter:extend URLToString
// goverter:extend MapDrivingLicense
type UserHTTPMapper interface {
	FromResponse(*domain.User) *entities.UserResponse
	FromResponseList([]*domain.User) []*entities.UserResponse
}

func UUIDToUUID(src uuid.UUID) string {
	return src.String()
}

func URLToString(src *url.URL) string {
	if src == nil {
		return ""
	}
	return src.String()
}
