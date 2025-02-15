package tree

import (
	"context"
	"testing"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestTreeRepository_Create(t *testing.T) {
	t.Run("should create a tree with default values", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)

		// when
		got, err := r.Create(context.Background(), func(tree *entities.Tree) (bool, error) {
			return true, nil
		})

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "", got.Species)
		assert.NotZero(t, got.ID)
		assert.WithinDuration(t, got.CreatedAt, time.Now(), time.Second)
		assert.WithinDuration(t, got.UpdatedAt, time.Now(), time.Second)
		assert.Nil(t, got.TreeCluster)
		assert.Nil(t, got.Sensor)
		assert.Equal(t, "", got.Number)
		assert.Equal(t, int32(0), got.PlantingYear)
		assert.Equal(t, float64(0), got.Latitude)
		assert.Equal(t, float64(0), got.Longitude)
		assert.Equal(t, "", got.Description)
		assert.Equal(t, false, got.Readonly)
		assert.Equal(t, entities.WateringStatusUnknown, got.WateringStatus)
	})

	t.Run("should create a tree with all values set, except images", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := TreeRepository{store: suite.Store, TreeMappers: mappers}
		sqlTreeCluster, clusterErr := suite.Store.GetTreeClusterByID(context.Background(), 1)
		if clusterErr != nil {
			t.Fatal(clusterErr)
		}

		treeCluster, err := mappers.tcMapper.FromSql(sqlTreeCluster)
		if err != nil {
			t.Fatal(err)
		}

		sensorID := "sensor-1"
		sqlSensor, sensorErr := suite.Store.GetSensorByID(context.Background(), sensorID)
		if sensorErr != nil {
			t.Fatal(sensorErr)
		}

		sensor, err := mappers.sMapper.FromSql(sqlSensor)
		if err != nil {
			t.Fatal(err)
		}

		// when
		got, err := r.Create(context.Background(), func(tree *entities.Tree) (bool, error) {
			tree.Species = "Oak"
			tree.Number = "T001"
			tree.PlantingYear = 2023
			tree.Latitude = 54.801539
			tree.Longitude = 9.446741
			tree.Description = "A newly planted oak tree"
			tree.WateringStatus = entities.WateringStatusGood
			tree.TreeCluster = treeCluster
			tree.Sensor = sensor
			return true, nil
		})

		treeClusterByTree, errClusterByTree := r.getTreeClusterByTreeID(context.Background(), got.ID)
		sensorByTree, errSensorByTree := r.GetSensorByTreeID(context.Background(), got.ID)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.NotZero(t, got.ID)
		assert.WithinDuration(t, got.CreatedAt, time.Now(), time.Second)
		assert.WithinDuration(t, got.UpdatedAt, time.Now(), time.Second)
		assert.NoError(t, errClusterByTree)
		assert.NotNil(t, treeClusterByTree)
		assert.Equal(t, treeCluster.ID, treeClusterByTree.ID)
		assert.NoError(t, errSensorByTree)
		assert.NotNil(t, sensorByTree)
		assert.Equal(t, sensor.ID, sensorByTree.ID)
		assert.Empty(t, got.Images)
		assert.Equal(t, "Oak", got.Species)
		assert.Equal(t, "T001", got.Number)
		assert.Equal(t, int32(2023), got.PlantingYear)
		assert.Equal(t, 54.801539, got.Latitude)
		assert.Equal(t, 9.446741, got.Longitude)
		assert.Equal(t, "A newly planted oak tree", got.Description)
		assert.Equal(t, false, got.Readonly)
		assert.Equal(t, entities.WateringStatusGood, got.WateringStatus)
	})

	t.Run("should return error if latitude is out of bounds", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)

		// when
		got, err := r.Create(context.Background(), func(tree *entities.Tree) (bool, error) {
			tree.Latitude = -200
			tree.Longitude = 0
			return true, nil
		})

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), storage.ErrInvalidLatitude.Error())
		assert.Nil(t, got)
	})

	t.Run("should return error if longitude is out of bounds", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)

		// when
		got, err := r.Create(context.Background(), func(tree *entities.Tree) (bool, error) {
			tree.Latitude = 0
			tree.Longitude = 200
			return true, nil
		})

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), storage.ErrInvalidLongitude.Error())
		assert.Nil(t, got)
	})

	t.Run("should return error if context is canceled", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.Create(ctx, func(tree *entities.Tree) (bool, error) {
			tree.Species = "Oak"
			return true, nil
		})

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should create tree and link images successfully", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := TreeRepository{store: suite.Store, TreeMappers: mappers}
		sqlImages, err := suite.Store.GetAllImages(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		images := mappers.iMapper.FromSqlList(sqlImages)

		sqlTreeCluster, clusterErr := suite.Store.GetTreeClusterByID(context.Background(), 1)
		if clusterErr != nil {
			t.Fatal(clusterErr)
		}
		treeCluster, err := mappers.tcMapper.FromSql(sqlTreeCluster)
		if err != nil {
			t.Fatal(err)
		}

		sensorID := "sensor-1"
		sqlSensor, sensorErr := suite.Store.GetSensorByID(context.Background(), sensorID)
		if sensorErr != nil {
			t.Fatal(sensorErr)
		}

		sensor, err := mappers.sMapper.FromSql(sqlSensor)
		if err != nil {
			t.Fatal(err)
		}

		// when
		tree, createErr := r.Create(context.Background(), func(tree *entities.Tree) (bool, error) {
			tree.Species = "Oak"
			tree.Number = "T001"
			tree.Latitude = 54.801539
			tree.Longitude = 9.446741
			tree.PlantingYear = 2023
			tree.Description = "Test tree with images"
			tree.TreeCluster = treeCluster
			tree.Sensor = sensor
			tree.Images = images
			return true, nil
		})

		treeClusterByTree, errClusterByTree := r.getTreeClusterByTreeID(context.Background(), tree.ID)
		sensorByTree, errSensorByTree := r.GetSensorByTreeID(context.Background(), tree.ID)

		// then
		assert.NoError(t, createErr)
		assert.NotNil(t, tree)
		assert.Equal(t, "Oak", tree.Species)
		assert.Equal(t, "T001", tree.Number)
		assert.Equal(t, 54.801539, tree.Latitude)
		assert.Equal(t, 9.446741, tree.Longitude)
		assert.Equal(t, int32(2023), tree.PlantingYear)
		assert.Equal(t, "Test tree with images", tree.Description)
		assert.NotEmpty(t, tree.Images)
		assert.NoError(t, errClusterByTree)
		assert.NotNil(t, treeClusterByTree)
		assert.Equal(t, treeCluster.ID, treeClusterByTree.ID)
		assert.NoError(t, errSensorByTree)
		assert.NotNil(t, sensorByTree)
		assert.Equal(t, sensor.ID, sensorByTree.ID)
		assert.Equal(t, false, tree.Readonly)
		for i, img := range tree.Images {
			assert.Equal(t, images[i].ID, img.ID)
			assert.Equal(t, images[i].URL, img.URL)
			assert.Equal(t, *images[i].Filename, *img.Filename)
			assert.Equal(t, *images[i].MimeType, *img.MimeType)
		}
	})
}
