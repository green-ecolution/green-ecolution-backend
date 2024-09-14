package sensor

import (
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	_ "github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:TimeToTime
type MqttHTTPMapper interface {
	ToResponse(src *domain.MqttPayload) *MqttPayloadResponse
	ToResponseList(src []*domain.MqttPayload) []*MqttPayloadResponse
	FromResponse(src *MqttPayloadResponse) *domain.MqttPayload
	FromResponseList(src []*MqttPayloadResponse) []*domain.MqttPayload
}
