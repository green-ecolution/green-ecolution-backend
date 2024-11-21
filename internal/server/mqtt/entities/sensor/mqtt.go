package sensor

type WatermarkResponse struct {
	Resistance int `json:"resistance"`
	Centibar   int `json:"centibar"`
	Depth      int `json:"depth"`
} // @Name Watermark

type MqttPayloadResponse struct {
	DeviceID    string              `json:"device_id"`
	Battery     float64             `json:"battery"`
	Humidity    int                 `json:"humidity"`
	Temperature int                 `json:"temperature"`
	Watermarks  []WatermarkResponse `json:"watermarks"`
} // @Name MqttPayload
