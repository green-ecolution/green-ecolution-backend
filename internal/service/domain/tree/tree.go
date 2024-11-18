package tree

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/tree"
)

type TreeService struct {
	treeRepo        storage.TreeRepository
	sensorRepo      storage.SensorRepository
	ImageRepo       storage.ImageRepository
	treeClusterRepo storage.TreeClusterRepository
	locator         service.GeoClusterLocator
	validator       *validator.Validate
}

func NewTreeService(
	repoTree storage.TreeRepository,
	repoSensor storage.SensorRepository,
	repoImage storage.ImageRepository,
	treeClusterRepo storage.TreeClusterRepository,
	geoClusterLocator service.GeoClusterLocator,
) service.TreeService {
	return &TreeService{

		treeRepo:        repoTree,
		sensorRepo:      repoSensor,
		ImageRepo:       repoImage,
		treeClusterRepo: treeClusterRepo,
		locator:         geoClusterLocator,
		validator:       validator.New(),
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

func (s *TreeService) GetBySensorID(ctx context.Context, id int32) (*entities.Tree, error) {
	tr, err := s.treeRepo.GetBySensorID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	return tr, nil
}

func (s *TreeService) Create(ctx context.Context, treeCreate *entities.TreeCreate) (*entities.Tree, error) {
	if err := s.validator.Struct(treeCreate); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}

	fn := make([]entities.EntityFunc[entities.Tree], 0)
	if treeCreate.TreeClusterID != nil {
		treeClusterID, err := s.treeClusterRepo.GetByID(ctx, *treeCreate.TreeClusterID)
		if err != nil {
			return nil, handleError(err)
		}
		fn = append(fn, tree.WithTreeCluster(treeClusterID))
	}

	if treeCreate.SensorID != nil {
		sensorID, err := s.sensorRepo.GetByID(ctx, *treeCreate.SensorID)
		if err != nil {
			return nil, handleError(err)
		}
		fn = append(fn, tree.WithSensor(sensorID))
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

	if treeCreate.TreeClusterID != nil {
		tc, err := s.treeClusterRepo.GetByID(ctx, *treeCreate.TreeClusterID)
		if err != nil {
			return nil, handleError(err)
		}
		if err = s.locator.UpdateCluster(ctx, tc); err != nil {
			return nil, handleError(err)
		}
	}
	return newTree, nil
}

func (s *TreeService) Delete(ctx context.Context, id int32) error {
	treeEntity, err := s.treeRepo.GetByID(ctx, id)
	if err != nil {
		return handleError(err)
	}
	if err := s.treeRepo.Delete(ctx, id); err != nil {
		return handleError(err)
	}
	if treeEntity.TreeCluster != nil {
		if err := s.locator.UpdateCluster(ctx, treeEntity.TreeCluster); err != nil {
			return handleError(err)
		}
	}
	return nil
}

func (s *TreeService) Update(ctx context.Context, id int32, tu *entities.TreeUpdate) (*entities.Tree, error) {
	if err := s.validator.Struct(tu); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}
	currentTree, err := s.treeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	// Check if the tree is readonly (imported from csv)
	// if currentTree.Readonly {
	// 	return nil, handleError(fmt.Errorf("tree with ID %d is readonly and cannot be updated", id))
	// }

	fn := make([]entities.EntityFunc[entities.Tree], 0)
	if tu.TreeClusterID != nil {
		var treeCluster *entities.TreeCluster
		treeCluster, err = s.treeClusterRepo.GetByID(ctx, *tu.TreeClusterID)
		if err != nil {
			return nil, handleError(fmt.Errorf("failed to find TreeCluster with ID %d: %w", *tu.TreeClusterID, err))
		}
		fn = append(fn, tree.WithTreeCluster(treeCluster))
	} else {
		fn = append(fn, tree.WithTreeCluster(nil))
	}

	if tu.SensorID != nil {
		var sensor *entities.Sensor
		sensor, err = s.sensorRepo.GetByID(ctx, *tu.SensorID)
		if err != nil {
			return nil, handleError(fmt.Errorf("failed to find Sensor with ID %d: %w", *tu.SensorID, err))
		}
		fn = append(fn, tree.WithSensor(sensor))
	} else {
		fn = append(fn, tree.WithSensor(nil))
	}

	fn = append(fn, tree.WithPlantingYear(tu.PlantingYear),
		tree.WithSpecies(tu.Species),
		tree.WithTreeNumber(tu.Number),
		tree.WithLatitude(tu.Latitude),
		tree.WithLongitude(tu.Longitude),
		tree.WithDescription(tu.Description))
	updatedTree, err := s.treeRepo.Update(ctx, id, fn...)
	if err != nil {
		return nil, handleError(err)
	}
	if currentTree.TreeCluster != nil {
		if err = s.locator.UpdateCluster(ctx, currentTree.TreeCluster); err != nil {
			return nil, handleError(err)
		}
	}
	if updatedTree.TreeCluster != nil {
		if err = s.locator.UpdateCluster(ctx, updatedTree.TreeCluster); err != nil {
			return nil, handleError(err)
		}
	}
	return updatedTree, nil
}

func handleError(err error) error {
	if errors.Is(err, storage.ErrEntityNotFound) {
		return service.NewError(service.NotFound, storage.ErrTreeNotFound.Error())
	}

	if errors.Is(err, storage.ErrSensorNotFound) {
		return service.NewError(service.NotFound, err.Error())
	}

	return service.NewError(service.InternalError, err.Error())
}

func (s *TreeService) Ready() bool {
	return s.treeRepo != nil && s.sensorRepo != nil
}
