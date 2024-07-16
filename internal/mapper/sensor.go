package mapper

import (
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities/sensor"
	response "github.com/green-ecolution/green-ecolution-backend/internal/service/entities/sensor"
	repo "github.com/green-ecolution/green-ecolution-backend/internal/storage/entities/sensor"
)

// goverter:converter
// goverter:extend TimeToTime
type MqttMapper interface {
	// goverter:autoMap Data
	FromEntity(src *repo.MqttEntity) *domain.MqttPayload
	FromEntityList(src []*repo.MqttEntity) []*domain.MqttPayload

	ToEntity(src *domain.MqttPayload) *repo.MqttPayloadEntity

	ToResponse(src *domain.MqttPayload) *response.MqttPayloadResponse
	ToResponseList(src []*domain.MqttPayload) []*response.MqttPayloadResponse
	FromResponse(src *response.MqttPayloadResponse) *domain.MqttPayload
}
