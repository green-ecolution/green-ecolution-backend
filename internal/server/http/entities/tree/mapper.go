package tree

import (
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities/tree"
	_ "github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:TimeToTime
type TreeHTTPMapper interface {
	ToResponse(src *domain.Tree) *TreeResponse
	ToResponseList(src []*domain.Tree) []*TreeResponse
	FromResponse(src *TreeResponse) *domain.Tree

	ToTreeSensorDataResponse(src *domain.TreeSensorData) *TreeSensorDataResponse
	ToTreeSensorDataResponseList(src []*domain.TreeSensorData) []*TreeSensorDataResponse

	ToTreeSensorPredictionResponse(src *domain.TreeSensorPrediction) *TreeSensorPredictionResponse
}
