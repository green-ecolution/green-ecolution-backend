package treecluster

import (
	"context"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestTreeClusterRepository_Update(t *testing.T) {
  t.Run("should update tree cluster", func(t *testing.T) {
    // given
    suite.ResetDB(t)
    suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
    r := NewTreeClusterRepository(suite.Store, mappers)

    // when
    got, err := r.Update(context.Background(), 1, WithName("updated"))

    // then
    assert.NoError(t, err)
    assert.NotNil(t, got)
    assert.Equal(t, "updated", got.Name)
  })

  t.Run("should return error when update tree cluster with non-existing id", func(t *testing.T) {
    // given
    r := NewTreeClusterRepository(suite.Store, mappers)

    // when
    got, err := r.Update(context.Background(), 99, WithName("updated"))

    // then
    assert.Error(t, err)
    assert.Nil(t, got)
  })

  t.Run("should return error when update tree cluster with negative id", func(t *testing.T) {
    // given
    r := NewTreeClusterRepository(suite.Store, mappers)

    // when
    got, err := r.Update(context.Background(), -1, WithName("updated"))

    // then
    assert.Error(t, err)
    assert.Nil(t, got)
  })

  t.Run("should return error if context is canceled", func(t *testing.T) {
    // given
    r := NewTreeClusterRepository(suite.Store, mappers)
    ctx, cancel := context.WithCancel(context.Background())
    cancel()

    // when
    got, err := r.Update(ctx, 1, WithName("updated"))

    // then
    assert.Error(t, err)
    assert.Nil(t, got)
  })

  t.Run("should not update tree cluster when no changes", func(t *testing.T) {
    // given
    suite.ResetDB(t)
    suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
    r := NewTreeClusterRepository(suite.Store, mappers)
    gotBefore, err := r.GetByID(context.Background(), 1)
    assert.NoError(t, err)

    // when
    got, err := r.Update(context.Background(), 1)

    // then
    assert.NoError(t, err)
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

    // when
    got, err := r.Update(context.Background(), 1, WithTrees(trees))

    // then
    assert.NoError(t, err)
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

    // when
    got, err := r.Update(context.Background(), 1, WithTrees(nil))

    // then
    assert.NoError(t, err)
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

    // when
    got, err := r.Update(context.Background(), 1, WithLatitude(utils.P(1.0)), WithLongitude(utils.P(1.0)))

    // then
    assert.NoError(t, err)
    assert.NotNil(t, got)
    assert.Equal(t, 1.0, got.Latitude)
    assert.Equal(t, 1.0, got.Longitude)
  })

  t.Run("should remove tree cluster coordinates", func(t *testing.T) {
    // given
    suite.ResetDB(t)
    suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
    r := NewTreeClusterRepository(suite.Store, mappers)

    // when
    got, err := r.Update(context.Background(), 1, WithLatitude(nil), WithLongitude(nil))

    // then
    assert.NoError(t, err)
    assert.NotNil(t, got)
    assert.Nil(t, got.Latitude)
    assert.Nil(t, got.Longitude)
  })
}
