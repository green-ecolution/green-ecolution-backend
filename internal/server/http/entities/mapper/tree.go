package mapper

import (
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:TimeToTime
// goverter:extend MapTreeClusterToID
type TreeHTTPMapper interface {
	// goverter:map TreeCluster TreeClusterID
	// goverter:ignore Sensor
	FromResponse(*domain.Tree) *entities.TreeResponse
	FromResponseList([]*domain.Tree) []*entities.TreeResponse
}

func MapTreeClusterToID(treeCluster *domain.TreeCluster) *int32 {
	if treeCluster == nil {
		return nil
	}
	return &treeCluster.ID
}
