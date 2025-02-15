package treecluster

import (
	"context"
	"sort"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
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

		got, totalCount, err := r.GetAll(ctx, "")

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

		ctx := context.WithValue(context.Background(), "page", int32(1))
		ctx = context.WithValue(ctx, "limit", int32(-1))

		got, totalCount, err := r.GetAll(ctx, "test-provider")

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.NotEmpty(t, got)
		assert.Equal(t, totalCount, int64(1))
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
		got, totalCount, err := r.GetAll(ctx, "")

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
		got, totalCount, err := r.GetAll(ctx, "")

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
		got, totalCount, err := r.GetAll(ctx, "")

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
		got, totalCount, err := r.GetAll(ctx, "")

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
		_, _, err := r.GetAll(ctx, "")

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

		if got.Latitude != nil {
			assert.Equal(t, allTestCluster[0].Latitude, *got.Latitude)
		} else {
			assert.Nil(t, got.Latitude)
		}

		if got.Longitude != nil {
			assert.Equal(t, allTestCluster[0].Longitude, *got.Longitude)
		} else {
			assert.Nil(t, got.Longitude)
		}

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
			assert.Equal(t, allTestCluster[i].Description, cluster.Description)

			if cluster.Latitude != nil {
				assert.Equal(t, allTestCluster[i].Latitude, *cluster.Latitude)
			} else {
				assert.Nil(t, cluster.Latitude)
			}

			if cluster.Longitude != nil {
				assert.Equal(t, allTestCluster[i].Longitude, *cluster.Longitude)
			} else {
				assert.Nil(t, cluster.Longitude)
			}

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
	Provider       string
	AdditionalInfo map[string]interface{}
}

var allTestCluster = []*testTreeCluster{
	{
		ID:             1,
		Name:           "Solitüde Strand",
		Address:        "Solitüde Strand",
		Description:    "Alle Bäume am Strand",
		MoistureLevel:  0.75,
		WateringStatus: entities.WateringStatusGood,
		Latitude:       54.82094,
		Longitude:      9.488152515892045,
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
		Latitude:       54.78805731048199,
		Longitude:      9.489022,
		SoilCondition:  entities.TreeSoilConditionSchluffig,
		RegionID:       1,
		TreeIDs:        []int32{4, 5, 6},
	},
	{
		ID:             3,
		Name:           "Flensburger Stadion",
		Address:        "Flensburger Stadion",
		Description:    "Alle Bäume in der Gegend des Stadions in Mürwik",
		MoistureLevel:  0.7,
		WateringStatus: entities.WateringStatusUnknown,
		Latitude:       54.802163,
		Longitude:      9.446398,
		SoilCondition:  entities.TreeSoilConditionSchluffig,
		RegionID:       1,
		TreeIDs:        []int32{16, 17, 18, 19, 20},
	},
	{
		ID:             4,
		Name:           "Campus Hochschule",
		Address:        "Thomas-Finke Straße",
		Description:    "Gruppe ist besonders anfällig",
		MoistureLevel:  0.1,
		WateringStatus: entities.WateringStatusGood,
		Latitude:       54.77578311851497,
		Longitude:      9.450294506300525,
		SoilCondition:  entities.TreeSoilConditionSchluffig,
		RegionID:       4,
		TreeIDs:        []int32{12, 13, 14, 15},
	},
	{
		ID:             5,
		Name:           "Mathildenstraße",
		Address:        "Mathildenstraße",
		Description:    "Sehr enge Straße und dadurch schlecht zu bewässern.",
		MoistureLevel:  0.4,
		WateringStatus: entities.WateringStatusBad,
		Latitude:       54.78219253876479,
		Longitude:      9.423978982828825,
		SoilCondition:  entities.TreeSoilConditionSchluffig,
		RegionID:       10,
		TreeIDs:        []int32{7, 8, 9, 10, 11},
	},
	{
		ID:             6,
		Name:           "Nordstadt",
		Address:        "Apenrader Straße",
		Description:    "Guter Baumbestand mit großen Kronen.",
		MoistureLevel:  0.6,
		WateringStatus: entities.WateringStatusUnknown,
		Latitude:       54.807162,
		Longitude:      9.423138,
		SoilCondition:  entities.TreeSoilConditionSandig,
		RegionID:       13,
		TreeIDs:        []int32{21, 22, 23, 24},
	},
	{
		ID:             7,
		Name:           "TSB Neustadt",
		Address:        "Ecknerstraße",
		Description:    "Kleiner Baumbestand.",
		MoistureLevel:  0.75,
		WateringStatus: entities.WateringStatusGood,
		Latitude:       54.797162,
		Longitude:      9.41962,
		SoilCondition:  entities.TreeSoilConditionSandig,
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
		Provider:       "test-provider",
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
