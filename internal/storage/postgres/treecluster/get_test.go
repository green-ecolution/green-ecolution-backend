package treecluster

import (
	"context"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestTreeClusterRepository_GetAll(t *testing.T) {
	t.Run("GetAll tree clusters should return tree clusters order by id", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewTreeClusterRepository(suite.Store, mappers)

		// when
		got, err := r.GetAll(context.Background())

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.NotEmpty(t, got)
		assert.Len(t, got, len(allTestCluster))
		for i, tc := range got {
			assert.Equal(t, allTestCluster[i].ID, tc.ID)
			assert.Equal(t, allTestCluster[i].Name, tc.Name)

			// assert region
			if allTestCluster[i].RegionID == -1 {
				assert.Nil(t, tc.Region)
				assert.NoError(t, err)
			} else {
				assert.NotNil(t, tc.Region)
				assert.Equal(t, allTestCluster[i].RegionID, tc.Region.ID)
			}

			// assert trees
			assert.Len(t, tc.Trees, len(allTestCluster[i].TreeIDs))
			if len(allTestCluster[i].TreeIDs) == 0 {
				assert.Empty(t, tc.Trees)
			}

			for j, tree := range tc.Trees {
				assert.NotZero(t, tree)
				assert.Equal(t, allTestCluster[i].TreeIDs[j], tree.ID)
			}
		}
	})

	t.Run("GetAll tree clusters with empty database should return empty list", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeClusterRepository(suite.Store, mappers)

		// when
		got, err := r.GetAll(context.Background())

		// then
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("GetAll tree clusters with context canceled exceeded should return error", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		_, err := r.GetAll(ctx)

		// then
		assert.Error(t, err)
	})
}

func TestTreeClusterRepository_GetByID(t *testing.T) {
	suite.ResetDB(t)
	suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
	t.Run("GetByID tree cluster should return tree cluster", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)

		// when
		got, err := r.GetByID(context.Background(), 1)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, allTestCluster[0].ID, got.ID)
		assert.Equal(t, allTestCluster[0].Name, got.Name)

		// assert region
		if allTestCluster[0].RegionID == -1 {
			assert.Nil(t, got.Region)
			assert.NoError(t, err)
		} else {
			assert.NotNil(t, got.Region)
			assert.Equal(t, allTestCluster[0].RegionID, got.Region.ID)
		}

		// assert trees
		assert.Len(t, got.Trees, len(allTestCluster[0].TreeIDs))
		if len(allTestCluster[0].TreeIDs) == 0 {
			assert.Empty(t, got.Trees)
		}

		for j, tree := range got.Trees {
			assert.NotZero(t, tree)
			assert.Equal(t, allTestCluster[0].TreeIDs[j], tree.ID)
		}
	})

	t.Run("GetByID tree cluster with non-existing id should return error", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)

		// when
		got, err := r.GetByID(context.Background(), 99)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("GetByID tree cluster with negative id should return error", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)

		// when
		got, err := r.GetByID(context.Background(), -1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("GetByID tree cluster with context canceled exceeded should return error", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetByID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestTreeClusterRepository_GetSensorByTreeClusterID(t *testing.T) {
	suite.ResetDB(t)
	suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
	t.Run("GetSensorByTreeClusterID should return sensor", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)

		// when
		got, err := r.GetSensorByTreeClusterID(context.Background(), 1)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, int32(1), got.ID)
		assert.Equal(t, entities.SensorStatusOnline, got.Status)
	})

	t.Run("GetSensorByTreeClusterID with non-existing id should return error", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)

		// when
		got, err := r.GetSensorByTreeClusterID(context.Background(), 99)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("GetSensorByTreeClusterID with negative id should return error", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)

		// when
		got, err := r.GetSensorByTreeClusterID(context.Background(), -1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("GetSensorByTreeClusterID with context canceled exceeded should return error", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetSensorByTreeClusterID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

type testTreeCluster struct {
	ID       int32
	Name     string
	RegionID int32
	TreeIDs  []int32
}

var allTestCluster = []*testTreeCluster{
	{
		ID:       1,
		Name:     "Solitüde Strand",
		RegionID: 1,
		TreeIDs:  []int32{1, 2, 3},
	},
	{
		ID:       2,
		Name:     "Sankt-Jürgen-Platz",
		RegionID: 1,
		TreeIDs:  []int32{4, 5, 6},
	},
	{
		ID:       3,
		Name:     "Flensburger Stadion",
		RegionID: 1,
		TreeIDs:  []int32{16, 17, 18, 19, 20},
	},
	{
		ID:       4,
		Name:     "Campus Hochschule",
		RegionID: 4,
		TreeIDs:  []int32{12, 13, 14, 15},
	},
	{
		ID:       5,
		Name:     "Mathildenstraße",
		RegionID: 10,
		TreeIDs:  []int32{7, 8, 9, 10, 11},
	},
	{
		ID:       6,
		Name:     "Nordstadt",
		RegionID: 13,
		TreeIDs:  []int32{21, 22, 23, 24},
	},
	{
		ID:       7,
		Name:     "TSB Neustadt",
		RegionID: 13,
	},
	{
		ID:       8,
		Name:     "Gewerbegebiet Süd",
		RegionID: -1, // no region
	},
}
