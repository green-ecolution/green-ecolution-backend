package mapper

import (
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:TimeToTime
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:UUIDToString
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:URLToString
type UserHTTPMapper interface {
	FromResponse(*domain.User) *entities.UserResponse
	FromResponseList([]*domain.User) []*entities.UserResponse
}
