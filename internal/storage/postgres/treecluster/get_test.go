package treecluster

import (
	"context"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreeClusterRepository_GetAll(t *testing.T) {
	t.Run("should return all tree clusters ordered by name without limitation", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewTreeClusterRepository(suite.Store, mappers)

		ctx := context.WithValue(context.Background(), "page", int32(1))
		ctx = context.WithValue(ctx, "limit", int32(-1))

		got, totalCount, err := r.GetAll(ctx)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.NotEmpty(t, got)
		assert.Len(t, got, len(allTestCluster))
		assert.Equal(t, totalCount, int64(len(allTestCluster)))

		sortedTestCluster := sortClusterByName(allTestCluster)

		for i, tc := range got {
			assert.Equal(t, sortedTestCluster[i].ID, tc.ID)
			assert.Equal(t, sortedTestCluster[i].Name, tc.Name)

			// assert region
			if sortedTestCluster[i].RegionID == -1 {
				assert.Nil(t, tc.Region)
				assert.NoError(t, err)
			} else {
				assert.NotNil(t, tc.Region)
				assert.Equal(t, sortedTestCluster[i].RegionID, tc.Region.ID)
			}

			// assert trees
			assert.Len(t, tc.Trees, len(sortedTestCluster[i].TreeIDs))
			if len(sortedTestCluster[i].TreeIDs) == 0 {
				assert.Empty(t, tc.Trees)
			}

			for j, tree := range tc.Trees {
				assert.NotZero(t, tree)
				assert.Equal(t, sortedTestCluster[i].TreeIDs[j], tree.ID)
			}
		}
	})

	t.Run("should return all tree clusters with provider", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewTreeClusterRepository(suite.Store, mappers)
		expectedCluster := allTestCluster[len(allTestCluster)-1]

		got, err := r.GetAllByProvider(context.Background(), "test-provider")

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.NotEmpty(t, got)
		assert.Len(t, got, 1)
		assert.Equal(t, expectedCluster.ID, got[0].ID)
		assert.Equal(t, expectedCluster.Name, got[0].Name)
		assert.Equal(t, expectedCluster.Provider, got[0].Provider)
		assert.Equal(t, expectedCluster.AdditionalInfo, got[0].AdditionalInfo)
	})

	t.Run("should return tree clusters ordered by name limited by 2 and with an offset of 2", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewTreeClusterRepository(suite.Store, mappers)

		ctx := context.WithValue(context.Background(), "page", int32(2))
		ctx = context.WithValue(ctx, "limit", int32(2))

		// when
		got, totalCount, err := r.GetAll(ctx)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.NotEmpty(t, got)
		assert.Len(t, got, 2)
		assert.Equal(t, totalCount, int64(len(allTestCluster)))

		sortedTestCluster := sortClusterByName(allTestCluster)[2:4]

		for i, tc := range got {
			assert.Equal(t, sortedTestCluster[i].ID, tc.ID)
			assert.Equal(t, sortedTestCluster[i].Name, tc.Name)
		}
	})

	t.Run("should return error on invalid page value", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewTreeClusterRepository(suite.Store, mappers)

		ctx := context.WithValue(context.Background(), "page", int32(0))
		ctx = context.WithValue(ctx, "limit", int32(2))

		// when
		got, totalCount, err := r.GetAll(ctx)

		// then
		assert.Error(t, err)
		assert.Empty(t, got)
		assert.Equal(t, totalCount, int64(0))
	})

	t.Run("should return error on invalid limit value", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewTreeClusterRepository(suite.Store, mappers)

		ctx := context.WithValue(context.Background(), "page", int32(2))
		ctx = context.WithValue(ctx, "limit", int32(0))

		// when
		got, totalCount, err := r.GetAll(ctx)

		// then
		assert.Error(t, err)
		assert.Empty(t, got)
		assert.Equal(t, totalCount, int64(0))
	})

	t.Run("should return empty slice when db is empty", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeClusterRepository(suite.Store, mappers)

		ctx := context.WithValue(context.Background(), "page", int32(2))
		ctx = context.WithValue(ctx, "limit", int32(2))

		// when
		got, totalCount, err := r.GetAll(ctx)

		// then
		assert.NoError(t, err)
		assert.Empty(t, got)
		assert.Equal(t, totalCount, int64(0))
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		_, _, err := r.GetAll(ctx)

		// then
		assert.Error(t, err)
	})
}

func TestTreeClusterRepository_GetByID(t *testing.T) {
	suite.ResetDB(t)
	suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
	t.Run("should return tree cluster by id", func(t *testing.T) {
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

	t.Run("should return error when tree cluster with non-existing id", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)

		// when
		got, err := r.GetByID(context.Background(), 99)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when tree cluster with negative id", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)

		// when
		got, err := r.GetByID(context.Background(), -1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when tree cluster with zero id", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)

		// when
		got, err := r.GetByID(context.Background(), 0)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
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

func TestTreeClusterRepository_GetByIDs(t *testing.T) {
	suite.ResetDB(t)
	suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")

	t.Run("should return tree clusters by ids", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)
		ids := []int32{1, 2}

		// when
		got, err := r.GetByIDs(context.Background(), ids)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Len(t, got, 2)

		for i, cluster := range got {
			assert.Equal(t, allTestCluster[i].ID, cluster.ID)
			assert.Equal(t, allTestCluster[i].Name, cluster.Name)

			// assert region
			if allTestCluster[i].RegionID == -1 {
				assert.Nil(t, cluster.Region)
				assert.NoError(t, err)
			} else {
				assert.NotNil(t, cluster.Region)
				assert.Equal(t, allTestCluster[i].RegionID, cluster.Region.ID)
			}

			// assert trees
			assert.Len(t, cluster.Trees, len(allTestCluster[i].TreeIDs))
			if len(allTestCluster[i].TreeIDs) == 0 {
				assert.Empty(t, cluster.Trees)
			}

			for j, tree := range cluster.Trees {
				assert.NotZero(t, tree)
				assert.Equal(t, allTestCluster[i].TreeIDs[j], tree.ID)
			}
		}
	})

	t.Run("should return empty list if no trees are found", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeClusterRepository(suite.Store, mappers)
		ids := []int32{99, 100, -1, 0}

		// when
		got, err := r.GetByIDs(context.Background(), ids)

		// then
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		trees, err := r.GetByIDs(ctx, []int32{1, 2})

		// then
		assert.Error(t, err)
		assert.Nil(t, trees)
	})
}

func TestTreeClusterRepository_GetAllLatestSensorDataByClusterID(t *testing.T) {
	suite.ResetDB(t)
	suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")

	t.Run("shold return all latest sensor data by cluster id", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)
		tcID := int32(50)

		// when
		got, err := r.GetAllLatestSensorDataByClusterID(context.Background(), tcID)

		// then
		assert.NoError(t, err)
		assert.Len(t, got, 2)
		assert.NotEqual(t, 34.0, got[0].Data.Battery) // based on seed
		assert.Equal(t, 99.0, got[0].Data.Battery)
		assert.NotEqual(t, 34.0, got[1].Data.Battery) // based on seed
		assert.Equal(t, 99.0, got[1].Data.Battery)
	})

	t.Run("shold return empty array when tree cluster not exists", func(t *testing.T) {
		// given
		r := NewTreeClusterRepository(suite.Store, mappers)
		tcID := int32(99)

		// when
		got, err := r.GetAllLatestSensorDataByClusterID(context.Background(), tcID)

		// then
		assert.NoError(t, err)
		assert.Empty(t, got)
	})
}

type testTreeCluster struct {
	ID             int32
	Name           string
	RegionID       int32
	TreeIDs        []int32
	Provider       string
	AdditionalInfo map[string]interface{}
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
	{
		ID:       50,
		Name:     "Gewerbegebiet Süd",
		RegionID: -1, // no region
		TreeIDs:  []int32{25, 26, 27, 28},
		Provider: "test-provider",
		AdditionalInfo: map[string]interface{}{
			"foo": "bar",
		},
	},
}

func sortClusterByName(data []*testTreeCluster) []*testTreeCluster {
	sorted := make([]*testTreeCluster, len(data))
	copy(sorted, data)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Name < sorted[j].Name
	})

	return sorted
}
