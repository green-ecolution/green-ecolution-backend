package region

import (
	"context"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestRegionRepository_GetAll(t *testing.T) {
	t.Run("Get all regions should return all regions ordered by id", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/region")
		r := NewRegionRepository(suite.Store, defaultRegionMappers())

		// when
		got, err := r.GetAll(context.Background())

		// then
		assert.NoError(t, err)
		assert.Equal(t, len(allTestRegions), len(got))
		for i, region := range got {
			assert.Equal(t, allTestRegions[i].ID, region.ID)
			assert.Equal(t, allTestRegions[i].Name, region.Name)
			assert.NotZero(t, region.CreatedAt)
			assert.NotZero(t, region.UpdatedAt)
		}
	})

	t.Run("Get all regions with empty database should return empty list", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewRegionRepository(suite.Store, defaultRegionMappers())

		// when
		got, err := r.GetAll(context.Background())

		// then
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("Get all regions with context canceled exceeded should return error", func(t *testing.T) {
		// given
		r := NewRegionRepository(suite.Store, defaultRegionMappers())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetAll(ctx)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestRegionRepository_GetByID(t *testing.T) {
	t.Run("Get region by id should return region", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/region")
		r := NewRegionRepository(suite.Store, defaultRegionMappers())
		shouldReturn := allTestRegions[0]

		// when
		got, err := r.GetByID(ctx, 1)

		// then
		assert.NoError(t, err)
		assert.Equal(t, shouldReturn.ID, got.ID)
		assert.Equal(t, shouldReturn.Name, got.Name)
		assert.NotZero(t, got.CreatedAt)
		assert.NotZero(t, got.UpdatedAt)
	})

	t.Run("Get region by id with non-existing id should return error", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewRegionRepository(suite.Store, defaultRegionMappers())

		// when
		got, err := r.GetByID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("Get region by id with negative id should return error", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewRegionRepository(suite.Store, defaultRegionMappers())

		// when
		got, err := r.GetByID(ctx, -1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("Get region by id with zero id should return error", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewRegionRepository(suite.Store, defaultRegionMappers())

		// when
		got, err := r.GetByID(ctx, 0)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("Get region by id with context canceled exceeded should return error", func(t *testing.T) {
		// given
		r := NewRegionRepository(suite.Store, defaultRegionMappers())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetByID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestRegionRepository_GetByName(t *testing.T) {
	tests := []struct {
		name string
		want *entities.Region
		args string
	}{
		{
			name: "Get region by name 'Mürwik' should return region",
			want: allTestRegions[0],
			args: "Mürwik",
		},
		{
			name: "Get region by name 'Fruerlund' should return region",
			want: allTestRegions[1],
			args: "Fruerlund",
		},
		{
			name: "Get region by name 'Jürgensby' should return region",
			want: allTestRegions[2],
			args: "Jürgensby",
		},
		{
			name: "Get region by name 'Sandberg' should return region",
			want: allTestRegions[3],
			args: "Sandberg",
		},
	}

	suite.ResetDB(t)
	suite.InsertSeed(t, "internal/storage/postgres/seed/test/region")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ctx := context.Background()
			r := NewRegionRepository(suite.Store, defaultRegionMappers())

			// when
			got, err := r.GetByName(ctx, tt.args)

			// then
			assert.NoError(t, err)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.NotZero(t, got.CreatedAt)
			assert.NotZero(t, got.UpdatedAt)
		})
	}

	t.Run("Get region by name with non-existing name should return error", func(t *testing.T) {
		// given
		ctx := context.Background()
		r := NewRegionRepository(suite.Store, defaultRegionMappers())

		// when
		got, err := r.GetByName(ctx, "Non-existing")

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("Get region by name with empty name should return error", func(t *testing.T) {
		// given
		ctx := context.Background()
		r := NewRegionRepository(suite.Store, defaultRegionMappers())

		// when
		got, err := r.GetByName(ctx, "")

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("Get region by name with context canceled exceeded should return error", func(t *testing.T) {
		// given
		r := NewRegionRepository(suite.Store, defaultRegionMappers())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetByName(ctx, "Mürwik")

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestRegionRepository_GetByPoint(t *testing.T) {
	t.Run("Get region by point should return region", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/region")
		r := NewRegionRepository(suite.Store, defaultRegionMappers())
		shouldReturn := allTestRegions[0]

		// when
		got, err := r.GetByPoint(ctx, 54.811925538974954, 9.484825422729664)

		// then
		assert.NoError(t, err)
		assert.Equal(t, shouldReturn.ID, got.ID)
		assert.Equal(t, shouldReturn.Name, got.Name)
		assert.NotZero(t, got.CreatedAt)
		assert.NotZero(t, got.UpdatedAt)
	})

	t.Run("Get region by point with non-existing point should return nil", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewRegionRepository(suite.Store, defaultRegionMappers())

		// when
		got, err := r.GetByPoint(ctx, 0, 0)

		// then
		assert.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("Get region by point with context canceled exceeded should return error", func(t *testing.T) {
		// given
		r := NewRegionRepository(suite.Store, defaultRegionMappers())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetByPoint(ctx, 54.413, 9.723)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

var allTestRegions = []*entities.Region{
	{
		ID:   1,
		Name: "Mürwik",
	},
	{
		ID:   2,
		Name: "Fruerlund",
	},
	{
		ID:   3,
		Name: "Jürgensby",
	},
	{
		ID:   4,
		Name: "Sandberg",
	},
	{
		ID:   5,
		Name: "Engelsby",
	},
	{
		ID:   6,
		Name: "Tarup",
	},
	{
		ID:   7,
		Name: "Altstadt",
	},
	{
		ID:   8,
		Name: "Südstadt",
	},
	{
		ID:   9,
		Name: "Weiche",
	},
	{
		ID:   10,
		Name: "Friesischer Berg",
	},
	{
		ID:   11,
		Name: "Westliche Höhe",
	},
	{
		ID:   12,
		Name: "Neustadt",
	},
	{
		ID:   13,
		Name: "Nordstadt",
	},
}
