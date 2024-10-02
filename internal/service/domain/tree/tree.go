package tree

import (
	"context"
	"errors"
	"fmt"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/tree"
	"log/slog"
)

type TreeService struct {
	treeRepo   storage.TreeRepository
	sensorRepo storage.SensorRepository
	ImageRepo  storage.ImageRepository
}

func NewTreeService(repoTree storage.TreeRepository, repoSensor storage.SensorRepository, repoImage storage.ImageRepository) service.TreeService {
	return &TreeService{
		treeRepo:   repoTree,
		sensorRepo: repoSensor,
		ImageRepo:  repoImage,
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
	var queue []*domain.Tree
	for i, csvTree := range trees {

		fmt.Printf("Coordinates hi: Latitude: %.6f, Longitude: %.6f\n", csvTree.Latitude, csvTree.Longitude)
		existingTree, err := s.treeRepo.GetByCoordinates(ctx, csvTree.Latitude, csvTree.Longitude)
		if err != nil || existingTree == nil {
			slog.Info("error finding tree by coordinates: %s", err)
			queue = append(queue, csvTree)
			slog.Info("Tree %d queued for later insertion", i+1)
			continue
		}

		if existingTree.Readonly {
			result := fmt.Sprintf("Tree %d (Latitude: %.6f, Longitude: %.6f) is read-only and cannot be overwritten. Please delete it manually.",
				i+1,
				csvTree.Latitude,
				csvTree.Longitude)
			slog.Error(result)
			return handleError(errors.New(result))
		}

		if existingTree.PlantingYear == csvTree.PlantingYear {
			err = s.updateTree(ctx, existingTree, csvTree, i)
		} else {
			err = s.deleteTree(ctx, existingTree, i)
			if existingTree.Sensor != nil {
				slog.Warn("Tree %d replaced, sensors are now unlinked.", i+1)
			}
			queue = append(queue, csvTree)
			slog.Info("Tree %d queued for creation", i+1)
		}
	}
	return s.processQueue(ctx, queue)
}

func (s *TreeService) updateTree(ctx context.Context, existingTree *domain.Tree, csvTree *domain.Tree, i int) error {
	//TODO: Update Attributes ?, the Attributes from scvRow should not be nil!!!!
	_, err := s.treeRepo.Update(ctx,
		existingTree.ID,
		tree.WithTreeNumber(csvTree.Number),

	)
	if err != nil {
		result := fmt.Sprintf("error updating tree %d: %v", i+1, err)
		return handleError(errors.New(result))
	}
	slog.Info("Tree %d updated with new attributes", i+1)

	return nil
}

func (s *TreeService) deleteTree(ctx context.Context, existingTree *domain.Tree, i int) error {
	err := s.treeRepo.DeleteAndUnlinkImages(ctx, existingTree.ID)
	if err != nil {
		result := fmt.Sprintf("error deleting tree %d: %v", i+1, err)
		return handleError(errors.New(result))
	}
	return nil
}

func (s *TreeService) processQueue(ctx context.Context, queue []*domain.Tree) error {
	if s.treeRepo == nil {
		slog.Error("treeRepo is nil")
		return errors.New("treeRepo is not initialized")
	}

	for _, newTree := range queue {
		if newTree == nil {
			slog.Error("newTree is nil")
			continue
		}
		fmt.Printf("Creating tree with Species: %s, Number: %v, Latitude: %.6f, Longitude: %.6f, PlantingYear: %d",
			newTree.Species, newTree.Number, newTree.Latitude, newTree.Longitude, newTree.PlantingYear)

		//TODO: the function repo.create returns a nil pointer dereference.
		_, err := s.treeRepo.Create(ctx,
			tree.WithSpecies(newTree.Species),
			tree.WithTreeNumber(newTree.Number),
			tree.WithLatitude(newTree.Latitude),
			tree.WithLongitude(newTree.Longitude),
			tree.WithPlantingYear(newTree.PlantingYear),
		)
		if err != nil {
			slog.Error("Error creating tree: %v", err)
			return err
		}
	}
	return nil
}
