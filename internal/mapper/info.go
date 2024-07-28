package mapper

import (
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities/info"
	response "github.com/green-ecolution/green-ecolution-backend/internal/service/entities/info"
	repo "github.com/green-ecolution/green-ecolution-backend/internal/storage/entities/info"
)

// goverter:converter
// goverter:extend TimeToTime URLToURL TimeDurationToTimeDuration StringToTime StringToURL StringToNetIP
// goverter:extend StringToDuration TimeToString NetURLToString NetIPToString TimeDurationToString
type InfoMapper interface {
	ToEntity(src *domain.App) *repo.AppEntity
	FromEntity(src *repo.AppEntity) *domain.App

	ToResponse(src *domain.App) *response.AppInfoResponse
	FromResponse(src *response.AppInfoResponse) *domain.App
}