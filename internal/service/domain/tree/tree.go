package tree

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/tree"
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
	t, err := s.treeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	return t, nil
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

func (s *TreeService) ImportTree(ctx context.Context, trees []*entities.Tree) error {
	var createQueue, updateQueue, deleteQueue []*entities.Tree

	for i, csvTree := range trees {
		// Log coordinates of the current tree
		slog.Info("Processing tree coordinates",
			"index", i+1,
			"latitude", csvTree.Latitude,
			"longitude", csvTree.Longitude)

		existingTree, err := s.treeRepo.GetByCoordinates(ctx, csvTree.Latitude, csvTree.Longitude)
		if err != nil || existingTree == nil {
			slog.Warn("Tree not found or error occurred",
				"index", i+1,
				"latitude", csvTree.Latitude,
				"longitude", csvTree.Longitude,
				"error", err)
			createQueue = append(createQueue, csvTree)
			slog.Info("Tree queued for creation", "index", i+1)
			continue
		}

		if existingTree.Readonly {
			slog.Error("Tree is read-only and cannot be overwritten",
				"index", i+1,
				"latitude", csvTree.Latitude,
				"longitude", csvTree.Longitude)
			return handleError(fmt.Errorf("tree %d (latitude: %.6f, longitude: %.6f) is read-only and cannot be overwritten. Please delete it manually",
				i+1, csvTree.Latitude, csvTree.Longitude))
		}

		if existingTree.PlantingYear == csvTree.PlantingYear {
			updateQueue = append(updateQueue, csvTree)
			slog.Info("Tree queued for update", "index", i+1)
		} else {
			deleteQueue = append(deleteQueue, existingTree)
			createQueue = append(createQueue, csvTree)
			slog.Info("Tree queued for deletion and recreation", "index", i+1)
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

func (s *TreeService) processDeleteQueue(ctx context.Context, deleteQueue []*entities.Tree) error {
	for i, treeToDelete := range deleteQueue {
		err := s.treeRepo.DeleteAndUnlinkImages(ctx, treeToDelete.ID)
		if err != nil {
			slog.Error("Error deleting tree",
				"index", i+1,
				"tree_id", treeToDelete.ID,
				"error", err)
			return handleError(fmt.Errorf("error deleting tree %d: %w", i+1, err))
		}
		if treeToDelete.Sensor != nil {
			slog.Warn("Tree is being replaced, sensors are now unlinked.",
				"index", i+1,
				"latitude", treeToDelete.Latitude,
				"longitude", treeToDelete.Longitude)
		}
	}
	return nil
}

func (s *TreeService) processUpdateQueue(ctx context.Context, updateQueue []*entities.Tree) error {
	for i, treeToUpdate := range updateQueue {
		_, err := s.treeRepo.Update(ctx, treeToUpdate.ID,
			tree.WithTreeNumber(treeToUpdate.Number),
			tree.WithSpecies(treeToUpdate.Species),
			tree.WithLatitude(treeToUpdate.Latitude),
			tree.WithLongitude(treeToUpdate.Longitude),
		)
		if err != nil {
			slog.Error("Error updating tree",
				"index", i+1,
				"tree_id", treeToUpdate.ID,
				"error", err)
			return handleError(err)
		}
		slog.Info("Tree updated successfully",
			"index", i+1,
			"tree_id", treeToUpdate.ID)
	}
	return nil
}

func (s *TreeService) processCreateQueue(ctx context.Context, createQueue []*entities.Tree) error {
	for i, newTree := range createQueue {
		if newTree == nil {
			slog.Error("newTree is nil", "index", i)
			continue
		}

		slog.Info("Creating tree",
			"Species", newTree.Species,
			"Number", newTree.Number,
			"Latitude", newTree.Latitude,
			"Longitude", newTree.Longitude,
			"PlantingYear", newTree.PlantingYear,
		)

		_, err := s.treeRepo.Create(ctx,
			tree.WithSpecies(newTree.Species),
			tree.WithTreeNumber(newTree.Number),
			tree.WithLatitude(newTree.Latitude),
			tree.WithLongitude(newTree.Longitude),
			tree.WithPlantingYear(newTree.PlantingYear),
		)

		if err != nil {
			slog.Error("Error creating tree", "index", i+1, "error", err)
			return handleError(fmt.Errorf("error creating tree %d: %w", i+1, err))
		}

		slog.Info("Tree created successfully", "index", i+1)
	}
	return nil
}
