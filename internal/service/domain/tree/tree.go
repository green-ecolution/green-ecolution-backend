package tree

import (
	"context"
	"errors"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	tree "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/tree"
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
	tree, err := s.treeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	return tree, nil
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
func (s *TreeService) Create(ctx context.Context, treeCreate *domain.TreeCreate) (*domain.Tree, error) {
	if treeCreate.PlantingYear == 0 {
		return nil, handleError(errors.New("plantingYear cannot be null or zero"))
	}
	if treeCreate.Species == "" {
		return nil, handleError(errors.New("species (Gattung) cannot be null or empty"))
	}
	if treeCreate.Number == "" {
		return nil, handleError(errors.New("tree Number (Baum Nr) cannot be null or empty"))
	}
	if treeCreate.Latitude == 0 || treeCreate.Longitude == 0 {
		return nil, handleError(errors.New("latitude and Longitude cannot be null or zero"))
	}

	fn := make([]domain.EntityFunc[domain.Tree], 0)
	if treeCreate.TreeClusterID != nil {
		treeClusterId, err := s.treeClusterRepo.GetByID(ctx, *treeCreate.TreeClusterID)
		if err != nil {
			return nil, handleError(err)
		}
		fn = append(fn, tree.WithTreeCluster(treeClusterId))
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
	return newTree, nil
}
