package tree

import (
	"context"
	"errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain/treecluster"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/tree"
)

type TreeService struct {
	treeRepo        storage.TreeRepository
	sensorRepo      storage.SensorRepository
	treeClusterRepo storage.TreeClusterRepository
	locator         *treecluster.GeoClusterLocator
}

func NewTreeService(
	repoTree storage.TreeRepository,
	repoSensor storage.SensorRepository,
	treeClusterRepo storage.TreeClusterRepository,
	geoClusterLocator *treecluster.GeoClusterLocator,
) service.TreeService {
	return &TreeService{
		treeRepo:        repoTree,
		sensorRepo:      repoSensor,
		treeClusterRepo: treeClusterRepo,
		locator:         geoClusterLocator,
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

func (s *TreeService) Create(ctx context.Context, treeCreate *entities.TreeCreate) (*entities.Tree, error) {
	if treeCreate.PlantingYear == 0 {
		return nil, handleError(errors.New("planting year cannot be null or zero"))
	}
	if treeCreate.Species == "" {
		return nil, handleError(errors.New("species cannot be null or empty"))
	}
	if treeCreate.Number == "" {
		return nil, handleError(errors.New("tree number cannot be null or empty"))
	}
	if treeCreate.Latitude == 0 || treeCreate.Longitude == 0 {
		return nil, handleError(errors.New("latitude and longitude cannot be null or zero"))
	}

	fn := make([]entities.EntityFunc[entities.Tree], 0)
	if treeCreate.TreeClusterID != nil {
		treeClusterID, err := s.treeClusterRepo.GetByID(ctx, *treeCreate.TreeClusterID)
		if err != nil {
			return nil, handleError(err)
		}
		fn = append(fn, tree.WithTreeCluster(treeClusterID))
	}
	fn = append(fn,
		tree.WithReadonly(treeCreate.Readonly),
		tree.WithPlantingYear(treeCreate.PlantingYear),
		tree.WithSpecies(treeCreate.Species),
		tree.WithTreeNumber(treeCreate.Number),
		tree.WithLatitude(treeCreate.Latitude),
		tree.WithLongitude(treeCreate.Longitude),
	)
	newTree, err := s.treeRepo.Create(ctx, fn...)
	if err != nil {
		return nil, handleError(err)
	}

	if err = s.locator.UpdateCluster(ctx, *treeCreate.TreeClusterID); err != nil {
		return nil, handleError(err)
	}

	return newTree, nil
}

func (s *TreeService) Delete(ctx context.Context, id int32) error {
	treeEntity, err := s.treeRepo.GetByID(ctx, id)
	if err != nil {
		return handleError(err)
	}

	treeClusterID := treeEntity.TreeCluster.ID
	if err := s.treeRepo.Delete(ctx, id); err != nil {
		return handleError(err)
	}

	return s.locator.UpdateCluster(ctx, treeClusterID)
}
