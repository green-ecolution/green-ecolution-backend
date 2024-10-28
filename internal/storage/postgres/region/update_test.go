package region

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegionRepository_Update(t *testing.T) {
	suite.ResetDB(t)
	suite.InsertSeed(t, "internal/storage/postgres/seed/test/region")

	t.Run("Update region should update region", func(t *testing.T) {
		// given
		r := NewRegionRepository(defaultFields.store, defaultFields.RegionMappers)

		// when
		got, err := r.Update(context.Background(), 1, WithName("test"))
		gotByID, _ := r.GetByID(context.Background(), 1)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "test", got.Name)
		assert.Equal(t, "test", gotByID.Name)
	})

	t.Run("Update region with empty name should return error", func(t *testing.T) {
		// given
		r := NewRegionRepository(defaultFields.store, defaultFields.RegionMappers)

		// when
		got, err := r.Update(context.Background(), 2, WithName(""))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("Update region with negative id should return error", func(t *testing.T) {
		// given
		r := NewRegionRepository(defaultFields.store, defaultFields.RegionMappers)

		// when
		got, err := r.Update(context.Background(), -1, WithName("test"))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("Update region without specifying name should not update", func(t *testing.T) {
		// given
		r := NewRegionRepository(defaultFields.store, defaultFields.RegionMappers)

		// when
		got, err := r.Update(context.Background(), 1)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "test", got.Name)
	})

	t.Run("Update region with non-existing id should return error", func(t *testing.T) {
		// given
		r := NewRegionRepository(defaultFields.store, defaultFields.RegionMappers)

		// when
		got, err := r.Update(context.Background(), 99, WithName("test"))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("Update region with context canceled exceeded should return error", func(t *testing.T) {
		// given
		r := NewRegionRepository(defaultFields.store, defaultFields.RegionMappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.Update(ctx, 1, WithName("test"))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}
