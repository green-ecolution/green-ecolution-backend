package entities

type Watermark struct {
	Resistance int `json:"resistance"`
	Centibar   int `json:"centibar"`
	Depth      int `json:"depth"`
}

type MqttPayload struct {
	DeviceID    string      `json:"device_id"`
	Battery     float64     `json:"battery"`
	Humidity    int         `json:"humidity"`
	Temperature int         `json:"temperature"`
	Watermarks  []Watermark `json:"watermarks"`
}
