package entities

type Watermark struct {
	Resistance int `json:"resistance"`
	Centibar   int `json:"centibar"`
	Depth      int `json:"depth"`
}

type MqttPayload struct {
	DeviceID    string      `json:"device_id"`
	Battery     float64     `json:"battery"`
	Humidity    float64     `json:"humidity"`
	Temperature float64     `json:"temperature"`
	Watermarks  []Watermark `json:"watermarks"`
}
