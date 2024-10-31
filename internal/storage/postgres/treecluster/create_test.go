package treecluster

import (
	"context"
	"testing"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestTreeClusterRepository_Create(t *testing.T) {
	t.Run("should create tree cluster with name", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewTreeClusterRepository(suite.Store, mappers)

		// when
		got, err := r.Create(context.Background(), WithName("test"))

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "test", got.Name)
		assert.NotZero(t, got.ID)
		assert.WithinDuration(t, got.CreatedAt, time.Now(), time.Second)
		assert.WithinDuration(t, got.UpdatedAt, time.Now(), time.Second)
		assert.Nil(t, got.Region)
		assert.Empty(t, got.Trees)
		assert.Equal(t, "", got.Address)
		assert.Equal(t, "", got.Description)
		assert.Equal(t, 0.0, got.MoistureLevel)
		assert.Nil(t, got.Latitude)
		assert.Nil(t, got.Longitude)
		assert.Equal(t, entities.WateringStatusUnknown, got.WateringStatus)
		assert.Equal(t, entities.TreeSoilConditionUnknown, got.SoilCondition)
		assert.False(t, got.Archived)
		assert.Nil(t, got.LastWatered)
	})

	t.Run("should create tree cluster with all values set", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewTreeClusterRepository(suite.Store, mappers)

		// when
		got, err := r.Create(context.Background(),
			WithName("test"),
			WithAddress("address"),
			WithDescription("description"),
			WithMoistureLevel(1.0),
			WithWateringStatus(entities.WateringStatusGood),
			WithSoilCondition(entities.TreeSoilConditionSchluffig),
		)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "test", got.Name)
		assert.NotZero(t, got.ID)
		assert.WithinDuration(t, got.CreatedAt, time.Now(), time.Second)
		assert.WithinDuration(t, got.UpdatedAt, time.Now(), time.Second)
		assert.Nil(t, got.Region)
		assert.Empty(t, got.Trees)
		assert.Equal(t, "address", got.Address)
		assert.Equal(t, "description", got.Description)
		assert.Equal(t, 1.0, got.MoistureLevel)
		assert.Nil(t, got.Latitude)
		assert.Nil(t, got.Longitude)
		assert.Equal(t, entities.WateringStatusGood, got.WateringStatus)
		assert.Equal(t, entities.TreeSoilConditionSchluffig, got.SoilCondition)
		assert.False(t, got.Archived)
		assert.Nil(t, got.LastWatered)
	})

	t.Run("should return tree cluster with trees and link tree cluster id to trees", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewTreeClusterRepository(suite.Store, mappers)
		testTrees, err := suite.Store.GetAllTrees(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		trees := mappers.treeMapper.FromSqlList(testTrees)[0:2]

		// when
		got, err := r.Create(context.Background(),
			WithName("test"),
			WithTrees(trees),
		)
		assert.NoError(t, err)

		sqlGotTrees, err := suite.Store.GetTreesByTreeClusterID(context.Background(), utils.P(got.ID))
		assert.NoError(t, err)

		assert.Len(t, sqlGotTrees, 2)
		assert.Equal(t, "test", got.Name)
		assert.NotZero(t, got.ID)
		assert.NotEmpty(t, got.Trees)

		for i, tree := range sqlGotTrees {
			assert.Equal(t, trees[i].ID, tree.ID)
			assert.NotNil(t, tree.TreeClusterID)
			assert.Equal(t, got.ID, *tree.TreeClusterID)
		}
	})

	t.Run("should return tree cluster with latitude and longitude", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewTreeClusterRepository(suite.Store, mappers)

		// when
		got, err := r.Create(context.Background(),
			WithName("test"),
			WithLatitude(utils.P(54.81269326939148)),
			WithLongitude(utils.P(9.484419532963013)),
		)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "test", got.Name)
		assert.NotZero(t, got.ID)
		assert.NotNil(t, got.Latitude)
		assert.NotNil(t, got.Longitude)
		assert.Equal(t, 54.81269326939148, *got.Latitude)
		assert.Equal(t, 9.484419532963013, *got.Longitude)
	})

	t.Run("should return error when tree cluster is invalid", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewTreeClusterRepository(suite.Store, mappers)

		// when
		got, err := r.Create(context.Background())

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when tree cluster has empty name", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewTreeClusterRepository(suite.Store, mappers)

		// when
		got, err := r.Create(context.Background(), WithName(""))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

  t.Run("should return error if context is canceled", func(t *testing.T) {
    // given
    suite.ResetDB(t)
    suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
    r := NewTreeClusterRepository(suite.Store, mappers)

    ctx, cancel := context.WithCancel(context.Background())
    cancel()

    // when
    got, err := r.Create(ctx, WithName("test"))

    // then
    assert.Error(t, err)
    assert.Nil(t, got)
  })
}

func TestTreeClusterRepository_LinkTreesToCluster(t *testing.T) {
	t.Run("should link trees to tree cluster", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewTreeClusterRepository(suite.Store, mappers)

		testTrees, err := suite.Store.GetAllTrees(context.Background())
		assert.NoError(t, err)
		trees := mappers.treeMapper.FromSqlList(testTrees)[0:2]

		tc, err := r.Create(context.Background(), WithName("test"))
		assert.NoError(t, err)

		// when
		err = r.LinkTreesToCluster(context.Background(), tc.ID, utils.Map(trees, func(t *entities.Tree) int32 {
			return t.ID
		}))
		assert.NoError(t, err)

		// then
		for _, tree := range testTrees[0:2] {
			// before
			if tree.TreeClusterID != nil {
				assert.NotEqual(t, tc.ID, *tree.TreeClusterID)
			}

			// after
			sqlTree, err := suite.Store.GetTreeByID(context.Background(), tree.ID)
			assert.NoError(t, err)
			assert.NotNil(t, sqlTree.TreeClusterID)
			assert.Equal(t, tc.ID, *sqlTree.TreeClusterID)
		}
	})

  t.Run("should return error when tree cluster not found", func(t *testing.T) {
    // given
    suite.ResetDB(t)
    suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
    r := NewTreeClusterRepository(suite.Store, mappers)

    testTrees, err := suite.Store.GetAllTrees(context.Background())
    assert.NoError(t, err)
    trees := mappers.treeMapper.FromSqlList(testTrees)[0:2]

    // when
    err = r.LinkTreesToCluster(context.Background(), 99, utils.Map(trees, func(t *entities.Tree) int32 {
      return t.ID
    }))

    // then
    assert.Error(t, err)
  })

  t.Run("should return error when context is canceled", func(t *testing.T) {
    // given
    suite.ResetDB(t)
    suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
    r := NewTreeClusterRepository(suite.Store, mappers)

    testTrees, err := suite.Store.GetAllTrees(context.Background())
    assert.NoError(t, err)
    trees := mappers.treeMapper.FromSqlList(testTrees)[0:2]

    tc, err := r.Create(context.Background(), WithName("test"))
    assert.NoError(t, err)

    ctx, cancel := context.WithCancel(context.Background())
    cancel()

    // when
    err = r.LinkTreesToCluster(ctx, tc.ID, utils.Map(trees, func(t *entities.Tree) int32 {
      return t.ID
    }))

    // then
    assert.Error(t, err)
  })
}
