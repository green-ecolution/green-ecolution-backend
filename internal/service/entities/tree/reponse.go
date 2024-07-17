package tree

import "github.com/green-ecolution/green-ecolution-backend/internal/service/entities/sensor"

type TreeSensorPredictionResponse struct {
	Tree             *TreeResponse                 `json:"tree,omitempty"`
	SensorPrediction *SensorPredictionResponse     `json:"sensor_prediction,omitempty"`
	SensorData       []*sensor.MqttPayloadResponse `json:"sensor_data,omitempty" extensions:"x-nullable"`
} //@Name TreeSensorPrediction

type TreeSensorDataResponse struct {
	Tree       *TreeResponse                 `json:"tree,omitempty"`
	SensorData []*sensor.MqttPayloadResponse `json:"sensor_data,omitempty" extensions:"x-nullable"`
} //@Name TreeSensorData
