package tree

import (
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

type TreeService struct {
	treeRepo   storage.TreeRepository
	sensorRepo storage.SensorRepository
}

func NewTreeService(repoTree storage.TreeRepository, repoSensor storage.SensorRepository) *TreeService {
	return &TreeService{
		treeRepo:   repoTree,
		sensorRepo: repoSensor,
	}
}

// func handleError(err error) error {
// 	if errors.Is(err, storage.ErrMongoDataNotFound) {
// 		return service.NewError(service.NotFound, err.Error())
// 	}
//
// 	return service.NewError(service.InternalError, err.Error())
// }

func (s *TreeService) Ready() bool {
	return s.treeRepo != nil && s.sensorRepo != nil
}
