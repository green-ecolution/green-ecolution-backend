package tree

import (
	"context"
	"errors"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

type TreeService struct {
	treeRepo   storage.TreeRepository
	sensorRepo storage.SensorRepository
}

func NewTreeService(repoTree storage.TreeRepository, repoSensor storage.SensorRepository) service.TreeService {
	return &TreeService{
		treeRepo:   repoTree,
		sensorRepo: repoSensor,
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

func (s *TreeService) ImportTree(ctx context.Context, trees []*domain.Tree) error {

	for i, tree := range trees {
		slog.Info("Tree %d: %+v", i+1, tree)

		slog.Info("Tree %d - Species: %s, TreeNumber: %d, Latitude: %.6f, Longitude: %.6f, PlantingYear: %d",
			i+1, tree.Species, tree.Number, tree.Latitude, tree.Longitude, tree.PlantingYear)

		//TODO: save the trees into database

		/*	_, err := s.treeRepo.Create(ctx,
				tree.WithSpecies(tree.Species),
				tree.WithTreeNumber(tree.Number),
				tree.WithLatitude(tree.Latitude),
				tree.WithLongitude(tree.Longitude),
				tree.WithPlantingYear(tree.PlantingYear),
			)
			if err != nil {
				slog.Info("Failed to create tree %d: %v", i+1, err)
				return err
			}*/
	}

	return nil
}
