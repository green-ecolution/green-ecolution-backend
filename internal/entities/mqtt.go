package entities

type Watermark struct {
	Centibar   int
	Resistance int
	Depth      int
}

type MqttPayload struct {
	DeviceID    string
	Battery     float64
	Humidity    float64
	Temperature float64
	Latitude    float64
	Longitude   float64
	Watermarks  []Watermark
}
