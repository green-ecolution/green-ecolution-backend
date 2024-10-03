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
	var createQueue, updateQueue, deleteQueue []*domain.Tree

	for i, csvTree := range trees {
		slog.Info("Coordinates: Latitude: %.6f, Longitude: %.6f\n", csvTree.Latitude, csvTree.Longitude)

		existingTree, err := s.treeRepo.GetByCoordinates(ctx, csvTree.Latitude, csvTree.Longitude)
		if err != nil || existingTree == nil {
			slog.Info("Error finding tree by coordinates: %v", err)
			createQueue = append(createQueue, csvTree)
			slog.Info("Tree %d queued for later creation", i+1)
			continue
		}

		if existingTree.Readonly {
			result := fmt.Sprintf("Tree %d (Latitude: %.6f, Longitude: %.6f) is read-only and cannot be overwritten. Please delete it manually.",
				i+1, csvTree.Latitude, csvTree.Longitude)
			slog.Error(result)
			return handleError(errors.New(result))
		}

		if existingTree.PlantingYear == csvTree.PlantingYear {
			updateQueue = append(updateQueue, csvTree)
			slog.Info("Tree %d queued for update", i+1)
		} else {
			deleteQueue = append(deleteQueue, existingTree)
			createQueue = append(createQueue, csvTree)
			slog.Info("Tree %d queued for deletion and recreation", i+1)
		}
	}

	if err := s.processDeleteQueue(ctx, deleteQueue); err != nil {
		return handleError(err)
	}
	if err := s.processUpdateQueue(ctx, updateQueue); err != nil {
		return handleError(err)
	}
	if err := s.processCreateQueue(ctx, createQueue); err != nil {
		return handleError(err)
	}

	return nil
}

func (s *TreeService) processDeleteQueue(ctx context.Context, deleteQueue []*domain.Tree) error {
	for i, treeToDelete := range deleteQueue {
		err := s.treeRepo.DeleteAndUnlinkImages(ctx, treeToDelete.ID)
		if err != nil {
			result := fmt.Sprintf("Error deleting tree %d: %v", i+1, err)
			slog.Error(result)
			return handleError(errors.New(result))
		}
		if treeToDelete.Sensor != nil {
			slog.Warn("Tree %d (Latitude: %.6f, Longitude: %.6f) is being replaced, sensors are now unlinked.",
				i+1, treeToDelete.Latitude, treeToDelete.Longitude)
		}
	}
	return nil
}
func (s *TreeService) processUpdateQueue(ctx context.Context, updateQueue []*domain.Tree) error {
	for i, treeToUpdate := range updateQueue {
		_, err := s.treeRepo.Update(ctx, treeToUpdate.ID,
			tree.WithTreeNumber(treeToUpdate.Number),
			tree.WithSpecies(treeToUpdate.Species),
			tree.WithLatitude(treeToUpdate.Latitude),
			tree.WithLongitude(treeToUpdate.Longitude),
		)
		if err != nil {
			result := fmt.Sprintf("Error updating tree %d: %v", i+1, err)
			slog.Error(result)
			return handleError(errors.New(result))
		}
		slog.Info("Tree %d updated successfully", i+1)
	}
	return nil
}
func (s *TreeService) processCreateQueue(ctx context.Context, createQueue []*domain.Tree) error {
	for i, newTree := range createQueue {
		if newTree == nil {
			slog.Error("newTree is nil")
			continue
		}
		info := fmt.Sprintf("Creating tree with Species: %s, Number: %v, Latitude: %.6f, Longitude: %.6f, PlantingYear: %d\n",
			newTree.Species, newTree.Number, newTree.Latitude, newTree.Longitude, newTree.PlantingYear)
		slog.Info(info)
		_, err := s.treeRepo.Create(ctx,
			tree.WithSpecies(newTree.Species),
			tree.WithTreeNumber(newTree.Number),
			tree.WithLatitude(newTree.Latitude),
			tree.WithLongitude(newTree.Longitude),
			tree.WithPlantingYear(newTree.PlantingYear),
		)
		if err != nil {
			result := fmt.Sprintf("Error creating tree %d: %v", i+1, err)
			slog.Error(result)
			return handleError(errors.New(result))
		}
		slog.Info("Tree %d created successfully", i+1)
	}
	return nil
}
