package treecluster

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
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/testutils"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/tree"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/treecluster/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"

	imgMapperImpl "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/image/mapper/generated"
	sensorMapperImpl "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor/mapper/generated"
	treeMapperImpl "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/tree/mapper/generated"
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
	Trees          []*randomTree                      `faker:"randomTrees"`
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

func initMappers() TreeClusterMappers {
	return NewTreeClusterRepositoryMappers(&generated.InternalTreeClusterRepoMapperImpl{}, &sensorMapperImpl.InternalSensorRepoMapperImpl{})
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

	mappers := initMappers()
	repo := NewTreeClusterRepository(str, mappers)

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

	imgMappers := image.NewImageRepositoryMappers(&imgMapperImpl.InternalImageRepoMapperImpl{})
	imgRepo := image.NewImageRepository(str, imgMappers)

	sensorRepo := sensor.NewSensorRepository(str, sensor.NewSensorRepositoryMappers(&sensorMapperImpl.InternalSensorRepoMapperImpl{}))

	// Create trees
	treeMappers := tree.NewTreeRepositoryMappers(&treeMapperImpl.InternalTreeRepoMapperImpl{}, &imgMapperImpl.InternalImageRepoMapperImpl{})
	treeRepo := tree.NewTreeRepository(str, treeMappers)
	for i, tree := range tc.Trees {
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

		// Create every second sensor
		var sensorID *int32
		sensorArg := &entities.CreateSensor{
			Status: tree.Sensor.Status,
		}
		if i%2 == 0 {
			sensorGot, err := sensorRepo.Create(context.Background(), sensorArg)
			assert.NoError(t, err)
			assert.NotNil(t, sensorGot)
			assert.NotZero(t, sensorGot.ID)
			assert.NotZero(t, sensorGot.CreatedAt)
			assert.NotZero(t, sensorGot.UpdatedAt)

			sensorID = &sensorGot.ID
		}

		arg := &entities.CreateTree{
			TreeClusterID:       got.ID,
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
		assertTree(t, tree, treeGot)
	}

	return got
}

func TestCreateTreeCluster(t *testing.T) {
	t.Parallel()
	t.Run("should create tree cluster", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			createTreeCluster(t, str)
		})
	})

	t.Run("should create tree cluster with no trees", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
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

			mappers := initMappers()
			repo := NewTreeClusterRepository(str, mappers)

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
		})
	})

	t.Run("should return error when creating tree cluster with invalid watering status", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			rtc := randomTreeCluster{}
			if err := faker.FakeData(&rtc); err != nil {
				t.Fatalf("error faking tree cluster data: %v", err)
			}

			rtc.WateringStatus = "invalid"
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

			mappers := initMappers()
			repo := NewTreeClusterRepository(str, mappers)

			_, err := repo.Create(context.Background(), &entities.CreateTreeCluster{
				WateringStatus: &tc.WateringStatus,
				MoistureLevel:  tc.MoistureLevel,
				Region:         tc.Region,
				Address:        tc.Address,
				Description:    tc.Description,
				Latitude:       tc.Latitude,
				Longitude:      tc.Longitude,
				SoilCondition:  &tc.SoilCondition,
			})
			assert.Error(t, err)
		})
	})

	t.Run("should return error when creating tree cluster with invalid soil condition", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			rtc := randomTreeCluster{}
			if err := faker.FakeData(&rtc); err != nil {
				t.Fatalf("error faking tree cluster data: %v", err)
			}

			rtc.SoilCondition = "invalid"
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

			mappers := initMappers()
			repo := NewTreeClusterRepository(str, mappers)

			_, err := repo.Create(context.Background(), &entities.CreateTreeCluster{
				WateringStatus: &tc.WateringStatus,
				MoistureLevel:  tc.MoistureLevel,
				Region:         tc.Region,
				Address:        tc.Address,
				Description:    tc.Description,
				Latitude:       tc.Latitude,
				Longitude:      tc.Longitude,
				SoilCondition:  &tc.SoilCondition,
			})
			assert.Error(t, err)
		})
	})

	t.Run("should return error when query failed", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
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

			mappers := initMappers()
			repo := NewTreeClusterRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.Create(context.Background(), &entities.CreateTreeCluster{
				WateringStatus: &tc.WateringStatus,
				MoistureLevel:  tc.MoistureLevel,
				Region:         tc.Region,
				Address:        tc.Address,
				Description:    tc.Description,
				Latitude:       tc.Latitude,
				Longitude:      tc.Longitude,
				SoilCondition:  &tc.SoilCondition,
			})
			assert.Error(t, err)
		})
	})
}

func TestGetAllTreeClusters(t *testing.T) {
	t.Parallel()
	t.Run("should get all tree clusters", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			tc := createTreeCluster(t, str)

			mappers := initMappers()
			repo := NewTreeClusterRepository(str, mappers)

			got, err := repo.GetAll(context.Background())
			assert.NoError(t, err)
			assert.NotEmpty(t, got)
			assertTreeCluster(t, tc, got[0])
		})
	})

	t.Run("should return empty list when no tree clusters found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMappers()
			repo := NewTreeClusterRepository(str, mappers)

			got, err := repo.GetAll(context.Background())
			assert.NoError(t, err)
			assert.Empty(t, got)
		})
	})

	t.Run("should return error when query failed", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)

			mappers := initMappers()
			repo := NewTreeClusterRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.GetAll(context.Background())
			assert.Error(t, err)
		})
	})
}

func TestGetTreeClusterByID(t *testing.T) {
	t.Parallel()
	t.Run("should get tree cluster by id", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			tc := createTreeCluster(t, str)

			mappers := initMappers()
			repo := NewTreeClusterRepository(str, mappers)

			got, err := repo.GetByID(context.Background(), tc.ID)
			assert.NoError(t, err)
			assert.NotNil(t, got)
			assertTreeCluster(t, tc, got)
		})
	})

	t.Run("should return error when tree cluster not found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMappers()
			repo := NewTreeClusterRepository(str, mappers)

			_, err := repo.GetByID(context.Background(), 999)
			assert.Error(t, err)
		})
	})

	t.Run("should return error when query failed", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			tc := createTreeCluster(t, str)

			mappers := initMappers()
			repo := NewTreeClusterRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.GetByID(context.Background(), tc.ID)
			assert.Error(t, err)
		})
	})
}

func TestGetSensorByTreeClusterID(t *testing.T) {
	t.Parallel()
	t.Run("should get sensor by tree cluster id", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			tc := createTreeCluster(t, str)

			for _, tree := range tc.Trees {
				if tree.Sensor == nil {
					continue
				}

				mappers := initMappers()
				repo := NewTreeClusterRepository(str, mappers)

				got, err := repo.GetSensorByTreeClusterID(context.Background(), tc.ID)
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assertSensor(t, tree.Sensor, got)
			}
		})
	})

	t.Run("should return error when sensor not found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMappers()
			repo := NewTreeClusterRepository(str, mappers)

			_, err := repo.GetSensorByTreeClusterID(context.Background(), 999)
			assert.Error(t, err)
		})
	})

	t.Run("should return error when query failed", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			tc := createTreeCluster(t, str)

			mappers := initMappers()
			repo := NewTreeClusterRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.GetSensorByTreeClusterID(context.Background(), tc.ID)
			assert.Error(t, err)
		})
	})
}

func TestUpdateTreeCluster(t *testing.T) {
	t.Parallel()
	t.Run("should update all tree cluster fields when all fields are provided", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			tc := createTreeCluster(t, str)

			mappers := initMappers()
			repo := NewTreeClusterRepository(str, mappers)
			args := &entities.UpdateTreeCluster{
				ID:             tc.ID,
				WateringStatus: utils.P(entities.TreeClusterWateringStatusGood),
				MoistureLevel:  utils.P(0.5),
				Region:         utils.P("new region"),
				Address:        utils.P("new address"),
				Description:    utils.P("new description"),
				Latitude:       utils.P(1.0),
				Longitude:      utils.P(1.0),
				SoilCondition:  utils.P(entities.TreeSoilConditionSchluffig),
			}
			want := &entities.TreeCluster{
				ID:             tc.ID,
				CreatedAt:      tc.CreatedAt,
				UpdatedAt:      time.Now(),
				WateringStatus: entities.TreeClusterWateringStatusGood,
				LastWatered:    tc.LastWatered,
				MoistureLevel:  0.5,
				Region:         "new region",
				Address:        "new address",
				Description:    "new description",
				Archived:       tc.Archived,
				Latitude:       1.0,
				Longitude:      1.0,
				SoilCondition:  entities.TreeSoilConditionSchluffig,
			}

			got, err := repo.Update(context.Background(), args)
			assert.NoError(t, err)
			assert.NotNil(t, got)
			assertTreeCluster(t, want, got)
		})
	})

	t.Run("should update only watering status when only watering status is provided", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			tc := createTreeCluster(t, str)

			mappers := initMappers()
			repo := NewTreeClusterRepository(str, mappers)
			args := &entities.UpdateTreeCluster{
				ID:             tc.ID,
				WateringStatus: utils.P(entities.TreeClusterWateringStatusGood),
				LastWatered:    nil,
				MoistureLevel:  nil,
				Region:         nil,
				Address:        nil,
				Description:    nil,
				Archived:       nil,
				Latitude:       nil,
				Longitude:      nil,
				SoilCondition:  nil,
			}
			want := &entities.TreeCluster{
				ID:             tc.ID,
				CreatedAt:      tc.CreatedAt,
				UpdatedAt:      time.Now(),
				WateringStatus: entities.TreeClusterWateringStatusGood,
				LastWatered:    tc.LastWatered,
				MoistureLevel:  tc.MoistureLevel,
				Region:         tc.Region,
				Address:        tc.Address,
				Description:    tc.Description,
				Archived:       tc.Archived,
				Latitude:       tc.Latitude,
				Longitude:      tc.Longitude,
				SoilCondition:  tc.SoilCondition,
			}

			got, err := repo.Update(context.Background(), args)
			assert.NoError(t, err)
			assert.NotNil(t, got)
			assertTreeCluster(t, want, got)
		})
	})

	t.Run("should update only moisture level when only moisture level is provided", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			tc := createTreeCluster(t, str)

			mappers := initMappers()
			repo := NewTreeClusterRepository(str, mappers)
			args := &entities.UpdateTreeCluster{
				ID:             tc.ID,
				MoistureLevel:  utils.P(0.5),
				LastWatered:    nil,
				WateringStatus: nil,
				Region:         nil,
				Address:        nil,
				Description:    nil,
				Archived:       nil,
				Latitude:       nil,
				Longitude:      nil,
				SoilCondition:  nil,
			}

			want := &entities.TreeCluster{
				ID:             tc.ID,
				CreatedAt:      tc.CreatedAt,
				UpdatedAt:      time.Now(),
				WateringStatus: tc.WateringStatus,
				LastWatered:    tc.LastWatered,
				MoistureLevel:  0.5,
				Region:         tc.Region,
				Address:        tc.Address,
				Description:    tc.Description,
				Archived:       tc.Archived,
				Latitude:       tc.Latitude,
				Longitude:      tc.Longitude,
				SoilCondition:  tc.SoilCondition,
			}

			got, err := repo.Update(context.Background(), args)
			assert.NoError(t, err)
			assert.NotNil(t, got)
			assertTreeCluster(t, want, got)
		})
	})

	t.Run("should update only provided fields", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			tc := createTreeCluster(t, str)

			mappers := initMappers()
			repo := NewTreeClusterRepository(str, mappers)
			args := &entities.UpdateTreeCluster{
				ID:             tc.ID,
				MoistureLevel:  utils.P(0.5),
				LastWatered:    nil,
				WateringStatus: nil,
				Region:         utils.P("new region"),
				Address:        nil,
				Description:    nil,
				Archived:       nil,
				Latitude:       nil,
				Longitude:      nil,
				SoilCondition:  utils.P(entities.TreeSoilConditionSchluffig),
			}

			want := &entities.TreeCluster{
				ID:             tc.ID,
				CreatedAt:      tc.CreatedAt,
				UpdatedAt:      time.Now(),
				WateringStatus: tc.WateringStatus,
				LastWatered:    tc.LastWatered,
				MoistureLevel:  0.5,
				Region:         "new region",
				Address:        tc.Address,
				Description:    tc.Description,
				Archived:       tc.Archived,
				Latitude:       tc.Latitude,
				Longitude:      tc.Longitude,
				SoilCondition:  entities.TreeSoilConditionSchluffig,
			}

			got, err := repo.Update(context.Background(), args)
			assert.NoError(t, err)
			assert.NotNil(t, got)
			assertTreeCluster(t, want, got)
		})
	})

	t.Run("should archive tree cluster if archive field is true", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			tc := createTreeCluster(t, str)

			mappers := initMappers()
			repo := NewTreeClusterRepository(str, mappers)
			args := &entities.UpdateTreeCluster{
				ID:             tc.ID,
				MoistureLevel:  utils.P(0.5),
				LastWatered:    nil,
				WateringStatus: utils.P(entities.TreeClusterWateringStatusBad),
				Region:         nil,
				Address:        nil,
				Description:    nil,
				Archived:       utils.P(true),
				Latitude:       nil,
				Longitude:      nil,
				SoilCondition:  nil,
			}

			want := &entities.TreeCluster{
				ID:             tc.ID,
				CreatedAt:      tc.CreatedAt,
				UpdatedAt:      time.Now(),
				WateringStatus: entities.TreeClusterWateringStatusBad,
				LastWatered:    tc.LastWatered,
				MoistureLevel:  0.5,
				Region:         tc.Region,
				Address:        tc.Address,
				Description:    tc.Description,
				Archived:       true,
				Latitude:       tc.Latitude,
				Longitude:      tc.Longitude,
				SoilCondition:  tc.SoilCondition,
			}

			got, err := repo.Update(context.Background(), args)
			assert.NoError(t, err)
			assert.NotNil(t, got)
			assertTreeCluster(t, want, got)
		})
	})

	t.Run("should return error when tree cluster not found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMappers()
			repo := NewTreeClusterRepository(str, mappers)
			args := &entities.UpdateTreeCluster{
				ID: 999,
			}

			_, err := repo.Update(context.Background(), args)
			assert.Error(t, err)
		})
	})

	t.Run("should return error when query failed", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			tc := createTreeCluster(t, str)

			mappers := initMappers()
			repo := NewTreeClusterRepository(str, mappers)
			args := &entities.UpdateTreeCluster{
				ID: tc.ID,
			}

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.Update(context.Background(), args)
			assert.Error(t, err)
		})
	})
}

func TestDeleteTreeCluster(t *testing.T) {
	t.Parallel()
	t.Run("should delete tree cluster by id", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			tc := createTreeCluster(t, str)

			mappers := initMappers()
			repo := NewTreeClusterRepository(str, mappers)

			err := repo.Delete(context.Background(), tc.ID)
			assert.NoError(t, err)

			_, err = repo.GetByID(context.Background(), tc.ID)
			assert.Error(t, err)
		})
	})

	t.Run("should return error when query failed", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			tc := createTreeCluster(t, str)

			mappers := initMappers()
			repo := NewTreeClusterRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			err = repo.Delete(context.Background(), tc.ID)
			assert.Error(t, err)
		})
	})
}

func TestArchiveTreeCluster(t *testing.T) {
  t.Parallel()
  t.Run("should archive tree cluster by id", func(t *testing.T) {
    testutils.WithTx(t, func(db *pgx.Conn) {
      str := createStore(db)
      tc := createTreeCluster(t, str)

      mappers := initMappers()
      repo := NewTreeClusterRepository(str, mappers)

      err := repo.Archive(context.Background(), tc.ID)
      assert.NoError(t, err)

      got, err := repo.GetByID(context.Background(), tc.ID)
      assert.NoError(t, err)
      assert.True(t, got.Archived)
    })
  })

  t.Run("should return error when query failed", func(t *testing.T) {
    testutils.WithTx(t, func(db *pgx.Conn) {
      str := createStore(db)
      tc := createTreeCluster(t, str)

      mappers := initMappers()
      repo := NewTreeClusterRepository(str, mappers)

      err := db.Close(context.Background())
      assert.NoError(t, err)

      err = repo.Archive(context.Background(), tc.ID)
      assert.Error(t, err)
    })
  })
}

func assertTreeCluster(t *testing.T, expected, actual *entities.TreeCluster) {
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

	assert.NotNil(t, actual.TreeCluster)
	assert.Equal(t, expected.TreeCluster.ID, actual.TreeCluster.ID)
	assert.Equal(t, expected.TreeCluster.CreatedAt, actual.TreeCluster.CreatedAt)
	assert.Equal(t, expected.TreeCluster.UpdatedAt, actual.TreeCluster.UpdatedAt)
	assert.Equal(t, expected.TreeCluster.WateringStatus, actual.TreeCluster.WateringStatus)
	assert.Equal(t, expected.TreeCluster.LastWatered, actual.TreeCluster.LastWatered)
	assert.Equal(t, expected.TreeCluster.MoistureLevel, actual.TreeCluster.MoistureLevel)
	assert.Equal(t, expected.TreeCluster.Region, actual.TreeCluster.Region)
	assert.Equal(t, expected.TreeCluster.Address, actual.TreeCluster.Address)
	assert.Equal(t, expected.TreeCluster.Description, actual.TreeCluster.Description)
	assert.Equal(t, expected.TreeCluster.Archived, actual.TreeCluster.Archived)
	assert.Equal(t, expected.TreeCluster.Latitude, actual.TreeCluster.Latitude)
	assert.Equal(t, expected.TreeCluster.Longitude, actual.TreeCluster.Longitude)
	assert.Equal(t, expected.TreeCluster.SoilCondition, actual.TreeCluster.SoilCondition)
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
