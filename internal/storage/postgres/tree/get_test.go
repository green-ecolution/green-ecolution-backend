package tree

import (
	"context"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/twpayne/go-geos"
)

func TestTreeRepository_GetAll(t *testing.T) {
	t.Run("should return all trees successfully", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)

		// when
		trees, err := r.GetAll(context.Background())

		// then
		assert.NoError(t, err)
		assert.NotNil(t, trees)
		assert.NotEmpty(t, trees)
		assert.Len(t, trees, len(testTrees))
		for i, expected := range trees {
			assertExpectedEqualToTree(t, expected, trees[i])
		}
	})

	t.Run("should return empty list when no trees exist", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)

		// when
		got, err := r.GetAll(context.Background())

		// then
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("should return error if context is canceled", func(t *testing.T) {
		// given
		r := NewTreeRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		trees, err := r.GetAll(ctx)

		// then
		assert.Error(t, err)
		assert.Nil(t, trees)
	})
}

func TestTreeRepository_GetByID(t *testing.T) {
	t.Run("should return the correct tree by ID", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)
		treeID := int32(1)

		sqlTreeCluster, clusterErr := suite.Store.GetTreeClusterByTreeID(context.Background(), treeID)
		if clusterErr != nil {
			t.Fatal(clusterErr)
		}
		treeCluster := mappers.tcMapper.FromSql(sqlTreeCluster)

		sqlSensor, sensorErr := suite.Store.GetSensorByTreeID(context.Background(), treeID)
		if sensorErr != nil {
			t.Fatal(sensorErr)
		}
		sensor := mappers.sMapper.FromSql(sqlSensor)

		// when
		tree, err := r.GetByID(context.Background(), 1)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, tree)
		assert.Equal(t, tree.TreeCluster.ID, treeCluster.ID)
		assert.Equal(t, tree.Sensor.ID, sensor.ID)
		assert.NotNil(t, tree.Sensor)
		assertExpectedEqualToTree(t, tree, testTrees[0])
	})

	t.Run("should return error if tree not found", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)

		// when
		tree, err := r.GetByID(context.Background(), 99)

		// then
		assert.Error(t, err)
		assert.Nil(t, tree)
	})

	t.Run("should return error if tree id is negative", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)

		// when
		tree, err := r.GetByID(context.Background(), -1)

		// then
		assert.Error(t, err)
		assert.Nil(t, tree)
	})

	t.Run("should return error if tree id is zero", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)

		// when
		tree, err := r.GetByID(context.Background(), 0)

		// then
		assert.Error(t, err)
		assert.Nil(t, tree)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewTreeRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		tree, err := r.GetByID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, tree)
	})
}

func TestTreeRepository_GetBySensorID(t *testing.T) {
	suite.ResetDB(t)
	suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")

	t.Run("should return the correct tree by linked sensor ID", func(t *testing.T) {
		// given
		r := NewTreeRepository(suite.Store, mappers)
		sensorID := "sensor-1"

		// when
		tree, err := r.GetBySensorID(context.Background(), sensorID)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, tree)
		assert.Equal(t, tree.Sensor.ID, sensorID)
		assertExpectedEqualToTree(t, tree, testTrees[0])
	})

	t.Run("should return error when sensor is not found", func(t *testing.T) {
		// given
		r := NewTreeRepository(suite.Store, mappers)
		sensorID := "sensor-notFound"

		// when
		tree, err := r.GetBySensorID(context.Background(), sensorID)

		// then
		assert.Error(t, err)
		assert.Nil(t, tree)
		assert.Equal(t, "sensor not found", err.Error())
	})

	t.Run("should return error when tree is not found", func(t *testing.T) {
		// given
		r := NewTreeRepository(suite.Store, mappers)
		sensorID := "sensor-4"

		// when
		tree, err := r.GetBySensorID(context.Background(), sensorID)

		// then
		assert.Error(t, err)
		assert.Nil(t, tree)
		assert.Equal(t, "entity not found", err.Error())
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewTreeRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		sensorID := "sensor-1"

		// when
		trees, err := r.GetBySensorID(ctx, sensorID)

		// then
		assert.Error(t, err)
		assert.Nil(t, trees)
	})
}

func TestTreeRepository_GetTreesByIDs(t *testing.T) {
	t.Run("should return trees successfully by IDs", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)
		ids := []int32{1, 2}
		expectedTrees := testTrees[:2]

		// when
		trees, err := r.GetTreesByIDs(context.Background(), ids)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, trees)
		assert.Len(t, trees, len(expectedTrees))
		for i, tree := range trees {
			assertExpectedEqualToTree(t, expectedTrees[i], tree)
		}
	})

	t.Run("should return empty list if no trees are found", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)
		ids := []int32{99, 100, -1, 0}

		// when
		trees, err := r.GetTreesByIDs(context.Background(), ids)

		// then
		assert.NoError(t, err)
		assert.Empty(t, trees)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewTreeRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		trees, err := r.GetTreesByIDs(ctx, []int32{1, 2})

		// then
		assert.Error(t, err)
		assert.Nil(t, trees)
	})
}

func TestTreeRepository_GetCenterPoint(t *testing.T) {
	t.Run("should successfully calculate and parse the center point", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)
		ids := []int32{1, 2, 3}
		centroids, err := suite.Store.CalculateGroupedCentroids(context.Background(), ids)
		if err != nil {
			t.Fatal(err)
		}
		g, errGeos := geos.NewGeomFromWKT(centroids)
		if errGeos != nil {
			t.Fatal(errGeos)
		}

		// when
		lat, long, errRepo := r.GetCenterPoint(context.Background(), ids)

		// then
		assert.NoError(t, errRepo)
		assert.Equal(t, g.X(), lat, "Latitude does not match")
		assert.Equal(t, g.Y(), long, "Longitude does not match")
	})

	t.Run("should return error if ids list is empty", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)

		// when
		lat, long, err := r.GetCenterPoint(context.Background(), []int32{})

		// then
		assert.Error(t, err)
		assert.Equal(t, 0.0, lat)
		assert.Equal(t, 0.0, long)
	})
}

func TestTreeRepository_GetByTreeClusterID(t *testing.T) {
	t.Run("should return trees successfully for a given tree cluster ID", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)
		treeClusterID := int32(1)
		expectedTrees := testTrees[:2]

		// when
		trees, err := r.GetByTreeClusterID(context.Background(), treeClusterID)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, trees)
		assert.NotEmpty(t, trees)
		assert.Len(t, trees, len(expectedTrees))
		for i, tree := range trees {
			assertExpectedEqualToTree(t, expectedTrees[i], tree)
		}
	})

	t.Run("should return error tree cluster ID is not found", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)

		// when
		trees, err := r.GetByTreeClusterID(context.Background(), 99)

		// then
		assert.Error(t, err)
		assert.Nil(t, trees)
	})

	t.Run("should return error tree cluster ID is zero", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)

		// when
		trees, err := r.GetByTreeClusterID(context.Background(), 0)

		// then
		assert.Error(t, err)
		assert.Nil(t, trees)
	})

	t.Run("should return error tree cluster ID is invalid", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)

		// when
		trees, err := r.GetByTreeClusterID(context.Background(), -1)

		// then
		assert.Error(t, err)
		assert.Nil(t, trees)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewTreeRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		trees, err := r.GetByTreeClusterID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, trees)
	})
}

func TestTreeRepository_GetByCoordinates(t *testing.T) {
	t.Run("should return tree successfully for given coordinates", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)
		testTree := testTrees[0]

		// when
		tree, err := r.GetByCoordinates(context.Background(), testTree.Latitude, testTree.Longitude)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, tree)
		assertExpectedEqualToTree(t, testTree, tree)
	})

	t.Run("should return error when no tree is found for given coordinates", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)
		latitude := 0.0
		longitude := 0.0

		// when
		tree, err := r.GetByCoordinates(context.Background(), latitude, longitude)

		// then
		assert.Error(t, err)
		assert.Nil(t, tree)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewTreeRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		tree, err := r.GetByCoordinates(ctx, 54.821248093376, 9.485710628517)

		// then
		assert.Error(t, err)
		assert.Nil(t, tree)
	})
}

func TestTreeRepository_GetAllImagesByID(t *testing.T) {
	t.Run("should return all images for the given tree ID", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)
		treeID := int32(1)

		// when
		images, err := r.GetAllImagesByID(context.Background(), treeID)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, images)
		assert.NotEmpty(t, images, "Images list should not be empty")
		for _, image := range images {
			assert.NotZero(t, image.ID, "Image ID should not be zero")
			assert.NotEmpty(t, image.URL, "Image URL should not be empty")
			assert.NotEmpty(t, image.Filename, "Image Filename should not be empty")
			assert.NotEmpty(t, image.MimeType, "Image MIME type should not be empty")
		}
	})

	t.Run("should return an empty list when the tree has no images", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)
		treeID := int32(2)

		// when
		images, err := r.GetAllImagesByID(context.Background(), treeID)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, images)
		assert.Empty(t, images, "Images list should be empty")
	})

	t.Run("should return an empty list when tree is not found", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)
		treeID := int32(999)

		// when
		images, err := r.GetAllImagesByID(context.Background(), treeID)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, images)
		assert.Empty(t, images, "Images list should be empty")
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewTreeRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		images, err := r.GetAllImagesByID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, images)
	})
}

func TestTreeRepository_GetSensorByTreeID(t *testing.T) {
	t.Run("should return sensor for the given tree ID", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)
		treeID := int32(1)

		// when
		sensor, err := r.GetSensorByTreeID(context.Background(), treeID)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, sensor, "Sensor should not be nil")
		assert.NotZero(t, sensor.ID, "Sensor ID should not be zero")
	})

	t.Run("should return ErrSensorNotFound when no sensor is linked to the tree", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)
		treeID := int32(2)

		// when
		sensor, err := r.GetSensorByTreeID(context.Background(), treeID)

		// then
		assert.ErrorIs(t, err, storage.ErrSensorNotFound, "Expected ErrSensorNotFound error")
		assert.Nil(t, sensor, "Sensor should be nil when not found")
	})

	t.Run("should return error when tree ID does not exist", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)
		invalidTreeID := int32(999)

		// when
		sensor, err := r.GetSensorByTreeID(context.Background(), invalidTreeID)

		// then
		assert.Error(t, err, "Expected an error when the tree ID does not exist")
		assert.Nil(t, sensor, "Sensor should be nil when tree ID does not exist")
	})

	t.Run("should return error when tree ID is negative", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)
		invalidTreeID := int32(-1)

		// when
		sensor, err := r.GetSensorByTreeID(context.Background(), invalidTreeID)

		// then
		assert.Error(t, err, "Expected an error when the tree ID is negative")
		assert.Nil(t, sensor, "Sensor should be nil when tree ID is negative")
	})

	t.Run("should return error when tree ID is zero", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)
		invalidTreeID := int32(0)

		// when
		sensor, err := r.GetSensorByTreeID(context.Background(), invalidTreeID)

		// then
		assert.Error(t, err, "Expected an error when the tree ID ID is zero")
		assert.Nil(t, sensor, "Sensor should be nil when tree ID ID is zero")
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewTreeRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		sensor, err := r.GetSensorByTreeID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, sensor)
	})
}

func TestTreeRepository_GetTreeClusterByTreeID(t *testing.T) {
	t.Run("should return tree cluster for the given tree ID", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := TreeRepository{store: suite.Store, TreeMappers: mappers}
		treeID := int32(1)

		// when
		treeCluster, err := r.getTreeClusterByTreeID(context.Background(), treeID)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, treeCluster, "TreeCluster should not be nil")
		assert.NotZero(t, treeCluster.ID, "TreeCluster ID should not be zero")
	})

	t.Run("should return ErrTreeClusterNotFound when no tree cluster is linked to the tree", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := TreeRepository{store: suite.Store, TreeMappers: mappers}
		treeID := int32(4)

		// when
		treeCluster, err := r.getTreeClusterByTreeID(context.Background(), treeID)

		// then
		assert.ErrorIs(t, err, storage.ErrTreeClusterNotFound, "Expected ErrTreeClusterNotFound error")
		assert.Nil(t, treeCluster, "TreeCluster should be nil when not found")
	})

	t.Run("should return error when tree ID does not exist", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := TreeRepository{store: suite.Store, TreeMappers: mappers}
		invalidTreeID := int32(999)

		// when
		treeCluster, err := r.getTreeClusterByTreeID(context.Background(), invalidTreeID)

		// then
		assert.Error(t, err, "Expected an error when the tree ID does not exist")
		assert.Nil(t, treeCluster, "TreeCluster should be nil when tree ID does not exist")
	})

	t.Run("should return error when tree ID is zero", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := TreeRepository{store: suite.Store, TreeMappers: mappers}
		invalidTreeID := int32(0)

		// when
		treeCluster, err := r.getTreeClusterByTreeID(context.Background(), invalidTreeID)

		// then
		assert.Error(t, err, "Expected an error when the tree ID is zero")
		assert.Nil(t, treeCluster, "TreeCluster should be nil when tree ID is zero")
	})

	t.Run("should return error when tree ID is negative", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := TreeRepository{store: suite.Store, TreeMappers: mappers}
		invalidTreeID := int32(-1)

		// when
		treeCluster, err := r.getTreeClusterByTreeID(context.Background(), invalidTreeID)

		// then
		assert.Error(t, err, "Expected an error when the tree ID is negative")
		assert.Nil(t, treeCluster, "TreeCluster should be nil when tree ID is negative")
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := TreeRepository{store: suite.Store, TreeMappers: mappers}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		treeCluster, err := r.getTreeClusterByTreeID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, treeCluster)
	})
}

func TestTreeRepository_FindNearestTree(t *testing.T) {
	t.Run("should return the nearest tree for given latitude and longitude", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)

		// when
		tree, err := r.GetByID(context.Background(), 2)
		assert.NoError(t, err)
		assert.NotNil(t, tree)

		sensorLatitude := 54.821517
		sensorLongitude := 9.487169

		// when
		nearestTree, errFind := r.FindNearestTree(context.Background(), sensorLatitude, sensorLongitude)

		// then
		assert.NoError(t, errFind, "Expected no error while finding the nearest tree")
		assert.NotNil(t, nearestTree, "Expected to find a nearest tree")
		assertExpectedEqualToTree(t, tree, nearestTree)
	})

	t.Run("should return error when no tree found within the required distance", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)

		sensorLatitude := 54.821535
		sensorLongitude := 9.487200

		// when
		nearestTree, err := r.FindNearestTree(context.Background(), sensorLatitude, sensorLongitude)

		// then
		assert.Error(t, err, "Expected error while finding the nearest tree")
		assert.Nil(t, nearestTree, "no tree should be found")
	})

	t.Run("should return error if context is canceled", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		tree, err := r.FindNearestTree(ctx, 54.82124518093376, 9.485702120628517)

		// then
		assert.Error(t, err, "Expected error when context is canceled")
		assert.Nil(t, tree, "Expected no tree to be returned when context is canceled")
	})
}

func assertExpectedEqualToTree(t *testing.T, expectedTree, tree *entities.Tree) {
	assert.Equal(t, expectedTree.ID, tree.ID, "ID does not match")
	assert.Equal(t, expectedTree.PlantingYear, tree.PlantingYear, "PlantingYear does not match")
	assert.Equal(t, expectedTree.Species, tree.Species, "Species does not match")
	assert.Equal(t, expectedTree.Number, tree.Number, "Number does not match")
	assert.Equal(t, expectedTree.Latitude, tree.Latitude, "Latitude does not match")
	assert.Equal(t, expectedTree.Longitude, tree.Longitude, "Longitude does not match")
	assert.Equal(t, expectedTree.Readonly, tree.Readonly, "Readonly does not match")
	assert.Equal(t, expectedTree.WateringStatus, tree.WateringStatus, "WateringStatus does not match")
	assert.Equal(t, expectedTree.Description, tree.Description, "Description does not match")
}

var testTrees = []*entities.Tree{
	{
		ID:             1,
		PlantingYear:   2021,
		Species:        "Quercus robur",
		Number:         "1005",
		Latitude:       54.82124518093376,
		Longitude:      9.485702120628517,
		Readonly:       true,
		WateringStatus: "unknown",
		Description:    "Sample description 1",
	},
	{
		ID:             2,
		PlantingYear:   2022,
		Species:        "Quercus robur",
		Number:         "1006",
		Latitude:       54.8215076622281,
		Longitude:      9.487153277881877,
		Readonly:       true,
		WateringStatus: "good",
		Description:    "Sample description 2",
	},
	{
		ID:             3,
		PlantingYear:   2023,
		Species:        "Betula pendula",
		Number:         "1010",
		Latitude:       54.78780993841013,
		Longitude:      9.444052105200551,
		Readonly:       false,
		WateringStatus: "bad",
		Description:    "Sample description 3",
	},
	{
		ID:             4,
		PlantingYear:   2020,
		Species:        "Quercus robur",
		Number:         "1008",
		Latitude:       54.787330993834613,
		Longitude:      9.4440523405200551,
		Readonly:       false,
		WateringStatus: "bad",
		Description:    "Sample description 4",
	},
}
