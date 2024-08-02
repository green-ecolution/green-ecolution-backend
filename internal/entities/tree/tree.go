package tree

import "github.com/green-ecolution/green-ecolution-backend/internal/entities/sensor"

type TreeLocation struct {
	Latitude       float64
	Longitude      float64
	Address        string
	AdditionalInfo string
}

type Tree struct {
	ID       string
	Species  string
	TreeNum  int
	Age      int
	Location TreeLocation
}

type TreeSensorData struct {
	Tree       *Tree
	SensorData []*sensor.MqttPayload
}

type TreeSensorPrediction struct {
	Tree             *Tree
	SensorPrediction *SensorPrediction
	SensorData       []*sensor.MqttPayload
}
