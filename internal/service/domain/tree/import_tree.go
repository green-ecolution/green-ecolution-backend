package tree

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
)

func (s *TreeService) ImportTree(ctx context.Context, trees []*entities.TreeImport) error {
	var deleteQueue []*entities.Tree
	var createQueue, updateQueue []*entities.TreeImport

	for i, csvTree := range trees {
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

		// if existingTree.Readonly {
		// 	slog.Error("Tree is read-only and cannot be overwritten",
		// 		"index", i+1,
		// 		"latitude", csvTree.Latitude,
		// 		"longitude", csvTree.Longitude)
		// 	return handleError(fmt.Errorf("tree %d (latitude: %.6f, longitude: %.6f) is read-only and cannot be overwritten. Please delete it manually",
		// 		i+1, csvTree.Latitude, csvTree.Longitude))
		// }

		if existingTree.PlantingYear == csvTree.PlantingYear {
			csvTree.TreeID = existingTree.ID
			updateQueue = append(updateQueue, csvTree)
			slog.Info("Tree queued for update", "index", i+1)
		} else {
			deleteQueue = append(deleteQueue, existingTree)
			createQueue = append(createQueue, csvTree)
			slog.Info("Tree queued for deletion and recreation", "index", i+1)
		}
	}

	if err := s.processDeleteQueue(ctx, deleteQueue); err != nil {
		return err
	}

	if err := s.processUpdateQueue(ctx, updateQueue); err != nil {
		return err
	}

	if err := s.processCreateQueue(ctx, createQueue); err != nil {
		return err
	}

	return nil
}

func (s *TreeService) processDeleteQueue(ctx context.Context, deleteQueue []*entities.Tree) error {
	for i, treeToDelete := range deleteQueue {
		err := s.Delete(ctx, treeToDelete.ID)
		if err != nil {
			slog.Error("Error deleting tree",
				"index", i+1,
				"tree_id", treeToDelete.ID,
				"error", err)
			return fmt.Errorf("error deleting tree %d: %w", i+1, err)
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

func (s *TreeService) processUpdateQueue(ctx context.Context, updateQueue []*entities.TreeImport) error {
	for i, treeToUpdate := range updateQueue {
		_, err := s.Update(ctx, treeToUpdate.TreeID, s.convertImportTreeToTreeUpdate(treeToUpdate))
		if err != nil {
			slog.Error("Error updating tree",
				"index", i+1,
				"tree_id", treeToUpdate.TreeID,
				"error", err)
			return err
		}
		slog.Info("Tree updated successfully",
			"index", i+1,
			"tree_id", treeToUpdate.TreeID)
	}
	return nil
}

func (s *TreeService) processCreateQueue(ctx context.Context, createQueue []*entities.TreeImport) error {
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
		_, err := s.Create(ctx, s.convertImportTreeToTreeCreate(newTree))
		if err != nil {
			slog.Error("Error creating tree", "index", i+1, "error", err)
			return fmt.Errorf("error creating tree %d: %w", i+1, err)
		}
		slog.Info("Tree created successfully", "index", i+1)
	}
	return nil
}

func (s *TreeService) convertImportTreeToTreeUpdate(tree *entities.TreeImport) *entities.TreeUpdate {
	treeUpdate := &entities.TreeUpdate{
		PlantingYear: tree.PlantingYear,
		Species:      tree.Species,
		Number:       tree.Number,
		Latitude:     tree.Latitude,
		Longitude:    tree.Longitude,
	}
	return treeUpdate
}

func (s *TreeService) convertImportTreeToTreeCreate(tree *entities.TreeImport) *entities.TreeCreate {
	treeCreate := &entities.TreeCreate{
		PlantingYear: tree.PlantingYear,
		Species:      tree.Species,
		Number:       tree.Number,
		Latitude:     tree.Latitude,
		Longitude:    tree.Longitude,
		Readonly:     true,
		Description:  "Dieser Baum wurde importiert",
	}
	return treeCreate
}
