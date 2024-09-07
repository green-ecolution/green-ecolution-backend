package tree

import (
	"context"
	"log/slog"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/image"
	mapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/testutils"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/treecluster"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

type randomTreeCluster struct {
	ID             int32                              `faker:"-"`
	CreatedAt      time.Time                          `faker:"-"`
	UpdatedAt      time.Time                          `faker:"-"`
	WateringStatus entities.TreeClusterWateringStatus `faker:"oneof: good, moderate, bad, unknown"`
	LastWatered    time.Time                          `faker:"-"`
	MoistureLevel  float64                            `faker:"oneof:0.1,0.2,0.3,0.4,0.5"`
	Region         string                             `faker:"word"`
	Address        string                             `faker:"oneof:address1,address2,address3,address4,address5"`
	Description    string                             `faker:"sentence"`
	Archived       bool                               `faker:"-"`
	Latitude       float64                            `faker:"lat"`
	Longitude      float64                            `faker:"long"`
	Trees          []*randomTree                      `faker:"-"`
	SoilCondition  entities.TreeSoilCondition         `faker:"oneof:schluffig"`
}

type randomTree struct {
	ID                  int32              `faker:"-"`
	CreatedAt           time.Time          `faker:"-"`
	UpdatedAt           time.Time          `faker:"-"`
	TreeCluster         *randomTreeCluster `faker:"-"`
	Sensor              *randomSensor      `faker:"randomSensor"`
	Images              []*randomImage     `faker:"randomImages"`
	Age                 int32              `faker:"oneof:1,2,3,4,5"`
	HeightAboveSeaLevel float64            `faker:"oneof:1.1,1.2,1.3,1.4,1.5"`
	PlantingYear        int32              `faker:"oneof:2020,2021,2022,2023,2024"`
	Species             string             `faker:"oneof:species1,species2,species3,species4,species5"`
	Number              int32              `faker:"oneof:1,2,3,4,5"`
	Latitude            float64            `faker:"lat"`
	Longitude           float64            `faker:"long"`
}

type randomSensor struct {
	ID        int32                 `faker:"-"`
	CreatedAt time.Time             `faker:"-"`
	UpdatedAt time.Time             `faker:"-"`
	Status    entities.SensorStatus `faker:"oneof:online,offline,unknown"`
}

type randomImage struct {
	ID        int32     `faker:"-"`
	CreatedAt time.Time `faker:"-"`
	UpdatedAt time.Time `faker:"-"`
	URL       string    `faker:"url"`
	Filename  *string   `faker:"word"`
	MimeType  *string   `faker:"oneof:image/png,image/jpeg"`
}

func initFaker() {
	faker.AddProvider("randomTrees", func(v reflect.Value) (interface{}, error) {
		trees := make([]*randomTree, 10)
		for i := 0; i < 10; i++ {
			tree := randomTree{}
			err := faker.FakeData(&tree)
			if err != nil {
				return nil, err
			}
			trees[i] = &tree
		}

		return trees, nil
	})

	faker.AddProvider("randomImages", func(v reflect.Value) (interface{}, error) {
		images := make([]*randomImage, 3)
		for i := 0; i < 3; i++ {
			img := randomImage{}
			err := faker.FakeData(&img)
			if err != nil {
				return nil, err
			}
			images[i] = &img
		}

		return images, nil
	})

	faker.AddProvider("randomSensor", func(v reflect.Value) (interface{}, error) {
		sensor := randomSensor{}
		err := faker.FakeData(&sensor)
		if err != nil {
			return nil, err
		}

		return &sensor, nil
	})
}

func TestMain(m *testing.M) {
	closeCon, _, err := testutils.SetupPostgresContainer()
	if err != nil {
		slog.Error("Error setting up postgres container", "error", err)
		os.Exit(1)
	}
	defer closeCon()
	initFaker()

	os.Exit(m.Run())
}

func createStore(db *pgx.Conn) *store.Store {
	return store.NewStore(db)
}

func initMappers() TreeMappers {
	return NewTreeRepositoryMappers(
		&mapper.InternalTreeRepoMapperImpl{},
		&mapper.InternalImageRepoMapperImpl{},
		&mapper.InternalSensorRepoMapperImpl{},
		&mapper.InternalTreeClusterRepoMapperImpl{},
	)
}

func treeClusterMappers() treecluster.TreeClusterMappers {
	return treecluster.NewTreeClusterRepositoryMappers(
		&mapper.InternalTreeClusterRepoMapperImpl{},
		&mapper.InternalSensorRepoMapperImpl{},
	)
}

func createTreeCluster(t *testing.T, str *store.Store) *entities.TreeCluster {
	rtc := randomTreeCluster{}
	if err := faker.FakeData(&rtc); err != nil {
		t.Fatalf("error faking tree cluster data: %v", err)
	}

	slog.Info("Creating tree cluster", "tc", rtc)
	tc := entities.TreeCluster{
		WateringStatus: rtc.WateringStatus,
		MoistureLevel:  rtc.MoistureLevel,
		Region:         rtc.Region,
		Address:        rtc.Address,
		Description:    rtc.Description,
		Latitude:       rtc.Latitude,
		Longitude:      rtc.Longitude,
		SoilCondition:  rtc.SoilCondition,
	}

	mapper := treeClusterMappers()
	repo := treecluster.NewTreeClusterRepository(str, mapper)

	got, err := repo.Create(context.Background(), &entities.CreateTreeCluster{
		WateringStatus: &tc.WateringStatus,
		MoistureLevel:  tc.MoistureLevel,
		Region:         tc.Region,
		Address:        tc.Address,
		Description:    tc.Description,
		Latitude:       tc.Latitude,
		Longitude:      tc.Longitude,
		SoilCondition:  &tc.SoilCondition,
	})
	assert.NoError(t, err)

	assert.NotNil(t, got)
	assert.NotZero(t, got.ID)

	assertTreeCluster(t, &tc, got)

	return got
}

func createTree(t *testing.T, str *store.Store, tc *entities.TreeCluster, hasSensor bool) *entities.Tree {
	rt := randomTree{}
	if err := faker.FakeData(&rt); err != nil {
		t.Fatalf("error faking tree data: %v", err)
	}

	mappers := initMappers()
	treeRepo := NewTreeRepository(str, mappers)

	imgMappers := image.NewImageRepositoryMappers(&mapper.InternalImageRepoMapperImpl{})
	imgRepo := image.NewImageRepository(str, imgMappers)

	sensorRepo := sensor.NewSensorRepository(str, sensor.NewSensorRepositoryMappers(&mapper.InternalSensorRepoMapperImpl{}))

	tree := entities.Tree{
		TreeCluster:         tc,
		Sensor:              &entities.Sensor{Status: entities.SensorStatusOnline},
		Images:              []*entities.Image{},
		Age:                 rt.Age,
		HeightAboveSeaLevel: rt.HeightAboveSeaLevel,
		PlantingYear:        rt.PlantingYear,
		Species:             rt.Species,
		Number:              rt.Number,
		Latitude:            rt.Latitude,
		Longitude:           rt.Longitude,
	}

	// Create Images
	for _, img := range tree.Images {
		arg := &entities.CreateImage{
			URL:      img.URL,
			Filename: img.Filename,
			MimeType: img.MimeType,
		}

		imgGot, err := imgRepo.Create(context.Background(), arg)
		assert.NoError(t, err)
		assert.NotNil(t, imgGot)
		assert.NotZero(t, imgGot.ID)
		assert.NotZero(t, imgGot.CreatedAt)
		assert.NotZero(t, imgGot.UpdatedAt)

		img.ID = imgGot.ID
		img.CreatedAt = imgGot.CreatedAt
		img.UpdatedAt = imgGot.UpdatedAt
	}

	imgIds := utils.Map(tree.Images, func(img *entities.Image) *int32 {
		return &img.ID
	})

	// Create Sensor
	var sensorID *int32
	if hasSensor {
		sensorArg := &entities.CreateSensor{
			Status: tree.Sensor.Status,
		}
		sensorGot, err := sensorRepo.Create(context.Background(), sensorArg)
		assert.NoError(t, err)
		assert.NotNil(t, sensorGot)
		assert.NotZero(t, sensorGot.ID)
		assert.NotZero(t, sensorGot.CreatedAt)
		assert.NotZero(t, sensorGot.UpdatedAt)

		sensorID = &sensorGot.ID
	}

	arg := &entities.CreateTree{
		TreeClusterID:       tc.ID,
		Age:                 tree.Age,
		HeightAboveSeaLevel: tree.HeightAboveSeaLevel,
		PlantingYear:        tree.PlantingYear,
		Species:             tree.Species,
		Latitude:            tree.Latitude,
		Longitude:           tree.Longitude,
		SensorID:            sensorID,
		ImageIDs:            imgIds,
	}

	treeGot, err := treeRepo.Create(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotNil(t, treeGot)
	assert.NotZero(t, treeGot.ID)

	assertTree(t, &tree, treeGot)

	return treeGot
}

func TestCreateTree(t *testing.T) {
	t.Parallel()
	t.Run("should create a tree", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			tc := createTreeCluster(t, str)
			createTree(t, str, tc, true)
		})
	})
}

func assertTreeCluster(t *testing.T, expected, actual *entities.TreeCluster) {
	if expected == nil {
		assert.Nil(t, expected)
		return
	}

	if actual == nil {
		assert.Nil(t, actual)
		return
	}

	assert.NotNil(t, actual)
	assert.NotZero(t, actual.ID)
	assert.NotZero(t, actual.CreatedAt)
	assert.NotZero(t, actual.UpdatedAt)

	assert.Equal(t, expected.WateringStatus, actual.WateringStatus)
	assert.Equal(t, expected.LastWatered, actual.LastWatered)
	assert.Equal(t, expected.MoistureLevel, actual.MoistureLevel)
	assert.Equal(t, expected.Region, actual.Region)
	assert.Equal(t, expected.Address, actual.Address)
	assert.Equal(t, expected.Description, actual.Description)
	assert.Equal(t, expected.Archived, actual.Archived)
	assert.Equal(t, expected.Latitude, actual.Latitude)
	assert.Equal(t, expected.Longitude, actual.Longitude)
	assert.Equal(t, expected.SoilCondition, actual.SoilCondition)
}

func assertTree(t *testing.T, expected, actual *entities.Tree) {
	assert.NotNil(t, actual)
	assert.NotZero(t, actual.ID)
	assert.NotZero(t, actual.CreatedAt)
	assert.NotZero(t, actual.UpdatedAt)

	if expected.TreeCluster != nil {
		assert.NotNil(t, actual.TreeCluster)
		assertTreeClusters(t, expected.TreeCluster, actual.TreeCluster)
	}

	if expected.Sensor != nil {
		assert.NotNil(t, actual.Sensor)
		assertSensor(t, expected.Sensor, actual.Sensor)
	}

	if expected.Images != nil && len(expected.Images) > 0 {
		assert.Len(t, actual.Images, len(expected.Images))
		for i := range expected.Images {
			assertImage(t, actual.Images[i], expected.Images[i])
		}
	}
}

func assertTreeClusters(t *testing.T, expected, actual *entities.TreeCluster) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.CreatedAt, actual.CreatedAt)
	assert.Equal(t, expected.UpdatedAt, actual.UpdatedAt)
	assert.Equal(t, expected.WateringStatus, actual.WateringStatus)
	assert.Equal(t, expected.LastWatered, actual.LastWatered)
	assert.Equal(t, expected.MoistureLevel, actual.MoistureLevel)
	assert.Equal(t, expected.Region, actual.Region)
	assert.Equal(t, expected.Address, actual.Address)
	assert.Equal(t, expected.Description, actual.Description)
	assert.Equal(t, expected.Archived, actual.Archived)
	assert.Equal(t, expected.Latitude, actual.Latitude)
	assert.Equal(t, expected.Longitude, actual.Longitude)
	assert.Equal(t, expected.SoilCondition, actual.SoilCondition)
}

func assertSensor(t *testing.T, got, want *entities.Sensor) {
	if want == nil {
		assert.Nil(t, got)
		return
	}

	if got == nil {
		assert.Nil(t, got)
		return
	}

	assert.NotZero(t, got.CreatedAt)
	assert.NotZero(t, got.UpdatedAt)
	assert.WithinDuration(t, got.CreatedAt, time.Now(), time.Second)
	assert.Equal(t, got.ID, want.ID)
	assert.Equal(t, got.Status, want.Status)
}

func assertImage(t *testing.T, got, want *entities.Image) {
	if got == nil {
		assert.Nil(t, got)
		return
	}

	if want == nil {
		assert.Nil(t, want)
		return
	}

	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.URL, got.URL)
	assert.Equal(t, want.Filename, got.Filename)
	assert.Equal(t, want.MimeType, got.MimeType)
}
