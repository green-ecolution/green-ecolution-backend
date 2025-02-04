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

		newRegion := &entities.Region{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "New Region",
		}
		now := time.Now()
		lat := 54.3
		long := 9.5

		tree1 := &entities.Tree{
			ID:             1,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			TreeCluster:    nil,
			Sensor:         nil,
			Images:         nil,
			Readonly:       false,
			PlantingYear:   2020,
			Species:        "Eiche",
			Number:         "1",
			Latitude:       54.31,
			Longitude:      9.51,
			WateringStatus: entities.WateringStatusGood,
			Description:    "Ein gesunder Baum",
		}

		tree2 := &entities.Tree{
			ID:             2,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			TreeCluster:    nil,
			Sensor:         nil,
			Images:         nil,
			Readonly:       false,
			PlantingYear:   2018,
			Species:        "Kiefer",
			Number:         "2",
			Latitude:       54.32,
			Longitude:      9.52,
			WateringStatus: entities.WateringStatusBad,
			Description:    "Braucht viel Wasser",
		}

		updateFn := func(tc *entities.TreeCluster) (bool, error) {
			tc.Name = "updated"
			tc.Address = "updated"
			tc.Description = "updated"
			tc.MoistureLevel = 4.2
			tc.WateringStatus = entities.WateringStatusBad
			tc.Archived = true
			tc.Region = newRegion
			tc.LastWatered = &now
			tc.Latitude = &lat
			tc.Longitude = &long
			tc.SoilCondition = entities.TreeSoilConditionLehmig
			tc.Trees = []*entities.Tree{tree1, tree2}
			return true, nil
		}

		// when
		updateErr := r.Update(context.Background(), 1, updateFn)
		got, getErr := r.GetByID(context.Background(), 1)

		// then
		assert.NoError(t, updateErr)
		assert.NoError(t, getErr)
		assert.NotNil(t, got)
		assert.Equal(t, "updated", got.Name)
		assert.Equal(t, "updated", got.Address)
		assert.Equal(t, "updated", got.Description)
		assert.Equal(t, 4.2, got.MoistureLevel)
		assert.Equal(t, entities.WateringStatusBad, got.WateringStatus)
		assert.Equal(t, true, got.Archived)
		assert.NotNil(t, got.Region)
		assert.Equal(t, newRegion.ID, got.Region.ID)
		assert.Equal(t, newRegion.Name, got.Region.Name)
		assert.NotNil(t, got.LastWatered)
		assert.WithinDuration(t, now, *got.LastWatered, time.Second)
		assert.NotNil(t, got.Latitude)
		assert.NotNil(t, got.Longitude)
		assert.Equal(t, lat, *got.Latitude)
		assert.Equal(t, long, *got.Longitude)
		assert.Equal(t, entities.TreeSoilConditionLehmig, got.SoilCondition)
		assert.NotNil(t, got.Trees)
		assert.Len(t, got.Trees, 2)
		assert.Equal(t, tree1.ID, got.Trees[0].ID)
		assert.Equal(t, tree2.ID, got.Trees[1].ID)
		assert.Equal(t, tree1.Species, got.Trees[0].Species)
		assert.Equal(t, tree2.Species, got.Trees[1].Species)
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
