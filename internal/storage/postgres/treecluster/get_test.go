package treecluster

import (
	"context"
	"sort"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestTreeClusterRepository_GetAll(t *testing.T) {
	t.Run("should return all tree clusters ordered by name", func(t *testing.T) {
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

	t.Run("should return empty slice when db is empty", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeClusterRepository(suite.Store, mappers)

		// when
		got, err := r.GetAll(context.Background())

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
		_, err := r.GetAll(ctx)

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
		assert.Equal(t, allTestCluster[0].Address, got.Address)
		assert.Equal(t, allTestCluster[0].Description, got.Description)
		assert.Equal(t, allTestCluster[0].MoistureLevel, got.MoistureLevel)
		assert.Equal(t, allTestCluster[0].WateringStatus, got.WateringStatus)
		assert.Equal(t, allTestCluster[0].SoilCondition, got.SoilCondition)
		assert.Equal(t, allTestCluster[0].Latitude, got.Latitude)
		assert.Equal(t, allTestCluster[0].Longitude, got.Longitude)

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
			assert.Equal(t, allTestCluster[i].Address, cluster.Address)
			assert.Equal(t, allTestCluster[i].MoistureLevel, cluster.MoistureLevel)
			assert.Equal(t, allTestCluster[i].WateringStatus, cluster.WateringStatus)
			assert.Equal(t, allTestCluster[i].SoilCondition, cluster.SoilCondition)
			assert.Equal(t, allTestCluster[i].Latitude, cluster.Latitude)
			assert.Equal(t, allTestCluster[i].Longitude, cluster.Longitude)
			assert.Equal(t, allTestCluster[i].Description, cluster.Description)

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
	Address        string
	Description    string
	MoistureLevel  float64
	WateringStatus entities.WateringStatus
	Latitude       float64
	Longitude      float64
	SoilCondition  entities.TreeSoilCondition
	RegionID       int32
	TreeIDs        []int32
}

var allTestCluster = []*testTreeCluster{
	{
		ID:             1,
		Name:           "Solitüde Strand",
		Address:        "Solitüde Strand",
		Description:    "Alle Bäume am Strand",
		MoistureLevel:  0.75,
		WateringStatus: entities.WateringStatusGood,
		Latitude:       0xc00019fdc0,
		Longitude:      0xc00019fdc8,
		SoilCondition:  entities.TreeSoilConditionSandig,
		RegionID:       1,
		TreeIDs:        []int32{1, 2, 3},
	},
	{
		ID:             2,
		Name:           "Sankt-Jürgen-Platz",
		Address:        "Ulmenstraße",
		Description:    "Bäume beim Sankt-Jürgen-Platz",
		MoistureLevel:  0.5,
		WateringStatus: entities.WateringStatusModerate,
		Latitude:       0xc0018fc3d0,
		Longitude:      0xc0018fc3d8,
		SoilCondition:  entities.TreeSoilConditionSchluffig,
		RegionID:       1,
		TreeIDs:        []int32{4, 5, 6},
	},
	{
		ID:             3,
		Name:           "Flensburger Stadion",
		Address:        "Address3",
		Description:    "Description3",
		MoistureLevel:  3.0,
		WateringStatus: entities.WateringStatusGood,
		Latitude:       0xc0017fc3d0,
		Longitude:      0xc0017fc3d8,
		SoilCondition:  entities.TreeSoilConditionSchluffig,
		RegionID:       1,
		TreeIDs:        []int32{16, 17, 18, 19, 20},
	},
	{
		ID:             4,
		Name:           "Campus Hochschule",
		Address:        "Address4",
		Description:    "Description4",
		MoistureLevel:  4.0,
		WateringStatus: entities.WateringStatusBad,
		Latitude:       0xc0016fc3d0,
		Longitude:      0xc0016fc3d8,
		SoilCondition:  entities.TreeSoilConditionLehmig,
		RegionID:       4,
		TreeIDs:        []int32{12, 13, 14, 15},
	},
	{
		ID:             5,
		Name:           "Mathildenstraße",
		Address:        "Address5",
		Description:    "Description5",
		MoistureLevel:  5.0,
		WateringStatus: entities.WateringStatusUnknown,
		Latitude:       0xc0015fc3d0,
		Longitude:      0xc0015fc3d8,
		SoilCondition:  entities.TreeSoilConditionSandig,
		RegionID:       10,
		TreeIDs:        []int32{7, 8, 9, 10, 11},
	},
	{
		ID:             6,
		Name:           "Nordstadt",
		Address:        "Address6",
		Description:    "Description6",
		MoistureLevel:  6.0,
		WateringStatus: entities.WateringStatusGood,
		Latitude:       0xc0014fc3d0,
		Longitude:      0xc0014fc3d8,
		SoilCondition:  entities.TreeSoilConditionLehmig,
		RegionID:       13,
		TreeIDs:        []int32{21, 22, 23, 24},
	},
	{
		ID:             7,
		Name:           "TSB Neustadt",
		Address:        "Address7",
		Description:    "Description7",
		MoistureLevel:  7.0,
		WateringStatus: entities.WateringStatusGood,
		Latitude:       0xc0013fc3d0,
		Longitude:      0xc0013fc3d8,
		SoilCondition:  entities.TreeSoilConditionSchluffig,
		RegionID:       13,
	},
	{
		ID:             8,
		Name:           "Gewerbegebiet Süd",
		Address:        "Address8",
		Description:    "Description8",
		MoistureLevel:  8.0,
		WateringStatus: entities.WateringStatusGood,
		Latitude:       0xc0012fc3d0,
		Longitude:      0xc0012fc3d8,
		SoilCondition:  entities.TreeSoilConditionLehmig,
		RegionID:       -1, // no region
	},
	{
		ID:             50,
		Name:           "Gewerbegebiet Süd",
		Address:        "Address9",
		Description:    "Description9",
		MoistureLevel:  9.0,
		WateringStatus: entities.WateringStatusGood,
		Latitude:       0xc0011fc3d0,
		Longitude:      0xc0011fc3d8,
		SoilCondition:  entities.TreeSoilConditionLehmig,
		RegionID:       -1, // no region
		TreeIDs:        []int32{25, 26, 27, 28},
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
