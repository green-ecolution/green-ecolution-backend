package tree

import (
	"context"
	"errors"
	"fmt"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/tree"
)

type TreeService struct {
	treeRepo        storage.TreeRepository
	sensorRepo      storage.SensorRepository
	treeClusterRepo storage.TreeClusterRepository
}

func NewTreeService(repoTree storage.TreeRepository, repoSensor storage.SensorRepository, treeClusterRepo storage.TreeClusterRepository) service.TreeService {
	return &TreeService{
		treeRepo:        repoTree,
		sensorRepo:      repoSensor,
		treeClusterRepo: treeClusterRepo,
	}
}

func (s *TreeService) GetAll(ctx context.Context) ([]*entities.Tree, error) {
	trees, err := s.treeRepo.GetAll(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	return trees, nil
}

func (s *TreeService) GetByID(ctx context.Context, id int32) (*entities.Tree, error) {
	tr, err := s.treeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	return tr, nil
}

func handleError(err error) error {
	if errors.Is(err, storage.ErrEntityNotFound) {
		return service.NewError(service.NotFound, err.Error())
	}

	return service.NewError(service.InternalError, err.Error())
}

func (s *TreeService) Ready() bool {
	return s.treeRepo != nil && s.sensorRepo != nil
}

func (s *TreeService) Delete(ctx context.Context, id int32) error {
	_, err := s.treeRepo.GetByID(ctx, id)
	if err != nil {
		return handleError(err)
	}
	err = s.treeRepo.DeleteAndUnlinkImages(ctx, id)
	if err != nil {
		return handleError(err)
	}
	return nil
}

func (s *TreeService) Update(ctx context.Context, id int32, tu *entities.TreeUpdate) (*entities.Tree, error) {
	currentTree, err := s.treeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}
	// Check if the tree is readonly (imported from csv)
	if currentTree.Readonly {
		return nil, handleError(fmt.Errorf("tree with ID %d is readonly and cannot be updated", id))
	}
	fn := make([]entities.EntityFunc[entities.Tree], 0)
	if tu.TreeClusterID != nil {
		treeCluster, err := s.treeClusterRepo.GetByID(ctx, *tu.TreeClusterID)
		if err != nil {
			return nil, handleError(fmt.Errorf("failed to find TreeCluster with ID %d: %w", *tu.TreeClusterID, err))
		}
		fn = append(fn, tree.WithTreeCluster(treeCluster))
	}
	if tu.PlantingYear != 0 {
		fn = append(fn, tree.WithPlantingYear(tu.PlantingYear))
	}
	if tu.Species != "" {
		fn = append(fn, tree.WithSpecies(tu.Species))
	}
	if tu.Number != "" {
		fn = append(fn, tree.WithTreeNumber(tu.Number))
	}
	if tu.Latitude != 0 && tu.Longitude != 0 {
		fn = append(fn, tree.WithLatitude(tu.Latitude), tree.WithLongitude(tu.Longitude))
	}

	updatedTree, err := s.treeRepo.Update(ctx, id, fn...)

	if err != nil {
		return nil, handleError(err)
	}
	return updatedTree, nil
}
