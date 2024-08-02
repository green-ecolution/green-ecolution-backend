package sensor

import (
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities/sensor"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:TimeToTime
type MqttMongoMapper interface {
	// goverter:autoMap Data
	FromEntity(src *MqttEntity) *domain.MqttPayload
	FromEntityList(src []*MqttEntity) []*domain.MqttPayload
	ToEntity(src *domain.MqttPayload) *MqttPayloadEntity
	ToEntityList(src []*domain.MqttPayload) []*MqttPayloadEntity
}
