package sensor

type SensorCreateRequest struct {
	Status SensorStatus `json:"status"`
	Type   string       `json:"type"`
} // @Name SensorCreateRequest

type SensorUpdateRequest struct {
  Status SensorStatus `json:"status"`
  Type   string       `json:"type"`
} // @Name SensorUpdateRequest

