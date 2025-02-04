package treecluster

import (
	"context"
	"testing"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestTreeClusterRepository_Update(t *testing.T) {
	t.Run("should update tree cluster", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewTreeClusterRepository(suite.Store, mappers)
		clusterID := int32(1)
		newName := "Updated Name"
		newAddress := "Updated Address"
		newDescription := "Updated Description"
		newMoistureLevel := 4.2
		newWateringStatus := entities.WateringStatusGood
		newSoilCondition := entities.TreeSoilConditionSchluffig
		newArchived := false
		newLastWatered := time.Now()
	
		// when
		updatedCluster, err := r.Update(
			context.Background(),
			clusterID,
			WithName(newName),
			WithAddress(newAddress),
			WithDescription(newDescription),
			WithMoistureLevel(newMoistureLevel),
			WithWateringStatus(newWateringStatus),
			WithSoilCondition(newSoilCondition),
			WithArchived(newArchived),
			WithLastWatered(newLastWatered),
		)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, updatedCluster)
		assert.Equal(t, newName, updatedCluster.Name, "Name should match")
		assert.Equal(t, newAddress, updatedCluster.Address, "Address should match")
		assert.Equal(t, newDescription, updatedCluster.Description, "Description should match")
		assert.Equal(t, newMoistureLevel, updatedCluster.MoistureLevel, "Moisture Level should match")
		assert.Equal(t, newWateringStatus, updatedCluster.WateringStatus, "Watering Status should match")
		assert.Equal(t, newSoilCondition, updatedCluster.SoilCondition, "Soil Condition should match")
		assert.Equal(t, newArchived, updatedCluster.Archived, "Archived status should match")
		assert.True(t, updatedCluster.LastWatered.Equal(newLastWatered), "Last Watered should match")
	})

	t.Run("should return error when update tree cluster with non-existing id", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)
		updateFn := func(tc *entities.TreeCluster) (bool, error) {
			tc.Name = "updated"
			return true, nil
		}

		// when
		err := r.Update(context.Background(), 99, updateFn)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when update tree cluster with negative id", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)
		updateFn := func(tc *entities.TreeCluster) (bool, error) {
			tc.Name = "updated"
			return true, nil
		}

		// when
		err := r.Update(context.Background(), -1, updateFn)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error if context is canceled", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)
		updateFn := func(tc *entities.TreeCluster) (bool, error) {
			tc.Name = "updated"
			return true, nil
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		err := r.Update(ctx, 1, updateFn)

		// then
		assert.Error(t, err)
	})

	t.Run("should not update tree cluster when no changes", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewTreeClusterRepository(suite.Store, mappers)
		gotBefore, err := r.GetByID(context.Background(), 1)
		assert.NoError(t, err)
		updateFn := func(tc *entities.TreeCluster) (bool, error) {
			return false, nil
		}

		// when
		updateErr := r.Update(context.Background(), 1, updateFn)
		got, getErr := r.GetByID(context.Background(), 1)

		// then
		assert.NoError(t, updateErr)
		assert.NoError(t, getErr)
		assert.NotNil(t, got)
		assert.Equal(t, gotBefore, got)
	})

	t.Run("should link trees to tree cluster", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewTreeClusterRepository(suite.Store, mappers)
		testTrees, err := suite.Store.GetAllTrees(context.Background())
		assert.NoError(t, err)
		trees := mappers.treeMapper.FromSqlList(testTrees)[0:2]
		updateFn := func(tc *entities.TreeCluster) (bool, error) {
			tc.Trees = trees
			return true, nil
		}

		// when
		updateErr := r.Update(context.Background(), 1, updateFn)
		got, getErr := r.GetByID(context.Background(), 1)

		// then
		assert.NoError(t, updateErr)
		assert.NoError(t, getErr)
		assert.NotNil(t, got)
		for _, tree := range testTrees[0:2] {
			assert.Equal(t, got.ID, *tree.TreeClusterID)
		}
	})

	t.Run("should unlink trees from tree cluster", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewTreeClusterRepository(suite.Store, mappers)
		beforeTreeCluster, err := r.GetByID(context.Background(), 1)
		assert.NoError(t, err)
		beforeTrees := beforeTreeCluster.Trees

		updateFn := func(tc *entities.TreeCluster) (bool, error) {
			tc.Trees = nil
			return true, nil
		}

		// when
		updateErr := r.Update(context.Background(), 1, updateFn)
		got, getErr := r.GetByID(context.Background(), 1)

		// then
		assert.NoError(t, updateErr)
		assert.NoError(t, getErr)
		assert.NotNil(t, got)
		for _, tree := range beforeTrees {
			actualTree, err := suite.Store.GetTreeByID(context.Background(), tree.ID)
			assert.NoError(t, err)
			assert.Nil(t, actualTree.TreeClusterID)
		}
	})

	t.Run("should update tree cluster coordinates", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewTreeClusterRepository(suite.Store, mappers)
		updateFn := func(tc *entities.TreeCluster) (bool, error) {
			tc.Latitude = utils.P(1.0)
			tc.Longitude = utils.P(1.0)
			return true, nil
		}

		// when
		updateErr := r.Update(context.Background(), 1, updateFn)
		got, getErr := r.GetByID(context.Background(), 1)

		// then
		assert.NoError(t, updateErr)
		assert.NoError(t, getErr)
		assert.NotNil(t, got)
		assert.NotNil(t, got.Latitude)
		assert.NotNil(t, got.Longitude)
		assert.Equal(t, 1.0, *got.Latitude)
		assert.Equal(t, 1.0, *got.Longitude)
	})

	t.Run("should remove tree cluster coordinates", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewTreeClusterRepository(suite.Store, mappers)
		updateFn := func(tc *entities.TreeCluster) (bool, error) {
			tc.Latitude = nil
			tc.Longitude = nil
			return true, nil
		}

		// when
		updateErr := r.Update(context.Background(), 1, updateFn)
		got, getErr := r.GetByID(context.Background(), 1)

		// then
		assert.NoError(t, updateErr)
		assert.NoError(t, getErr)
		assert.NotNil(t, got)
		assert.Nil(t, got.Latitude)
		assert.Nil(t, got.Longitude)
	})

	t.Run("should return error when updateFn is nil", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)

		// when
		err := r.Update(context.Background(), 1, nil)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when updateFn returns error", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)
		updateFn := func(tc *entities.TreeCluster) (bool, error) {
			return true, assert.AnError
		}

		// when
		err := r.Update(context.Background(), 1, updateFn)

		// then
		assert.Error(t, err)
	})

	t.Run("should not update when updateFn returns false", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)
		updateFn := func(tc *entities.TreeCluster) (bool, error) {
			return false, nil
		}

		// when
		updateErr := r.Update(context.Background(), 1, updateFn)
		got, getErr := r.GetByID(context.Background(), 1)

		// then
		assert.NoError(t, updateErr)
		assert.NoError(t, getErr)
		assert.NotNil(t, got)
	})

	t.Run("should not rollback when updateFn returns false", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewTreeClusterRepository(suite.Store, mappers)
		updateFn := func(tc *entities.TreeCluster) (bool, error) {
			tc.Name = "updated"
			return false, nil
		}

		// when
		err := r.Update(context.Background(), 1, updateFn)
		got, getErr := r.GetByID(context.Background(), 1)

		// then
		assert.NoError(t, err)
		assert.NoError(t, getErr)
		assert.NotNil(t, got)
		assert.NotEqual(t, "updated", got.Name)
	})
}
