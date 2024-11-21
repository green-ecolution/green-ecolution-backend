package entities

type Watermark struct {
	Centibar   int
	Resistance int
	Depth      int
}

type MqttPayload struct {
	DeviceID    string
	Battery     float64
	Humidity    int
	Temperature int
	Watermarks  []Watermark
}
