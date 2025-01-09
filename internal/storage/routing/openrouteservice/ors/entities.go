package ors

type OrsDirectionRequest struct {
	Coordinates [][]float64 `json:"coordinates"`
	Units       string      `json:"units"`
	Language    string      `json:"language"`
}
