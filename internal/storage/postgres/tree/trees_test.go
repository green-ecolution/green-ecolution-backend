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

	// Map random trees to entities.Tree
	for _, tree := range rtc.Trees {
		sensor := entities.Sensor{
			Status: tree.Sensor.Status,
		}

		images := utils.Map(tree.Images, func(img *randomImage) *entities.Image {
			return &entities.Image{
				URL:      img.URL,
				Filename: img.Filename,
				MimeType: img.MimeType,
			}
		})

		t := entities.Tree{
			TreeCluster:         got,
			Age:                 tree.Age,
			HeightAboveSeaLevel: tree.HeightAboveSeaLevel,
			PlantingYear:        tree.PlantingYear,
			Species:             tree.Species,
			Latitude:            tree.Latitude,
			Longitude:           tree.Longitude,
			Sensor:              &sensor,
			Images:              images,
		}
		got.Trees = append(got.Trees, &t)
	}

	return got
}

func createTrees(t *testing.T, str *store.Store, tc *entities.TreeCluster) {
	mappers := initMappers()
	treeRepo := NewTreeRepository(str, mappers)

	imgMappers := image.NewImageRepositoryMappers(&mapper.InternalImageRepoMapperImpl{})
	imgRepo := image.NewImageRepository(str, imgMappers)

	sensorRepo := sensor.NewSensorRepository(str, sensor.NewSensorRepositoryMappers(&mapper.InternalSensorRepoMapperImpl{}))

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

		// Create Sensor for every 2nd tree
		var sensor *entities.Sensor
		if i%2 == 0 {
			sensorArg := &entities.CreateSensor{
				Status: tree.Sensor.Status,
			}
			sensorGot, err := sensorRepo.Create(context.Background(), sensorArg)
			assert.NoError(t, err)
			assert.NotNil(t, sensorGot)
			assert.NotZero(t, sensorGot.ID)
			assert.NotZero(t, sensorGot.CreatedAt)
			assert.NotZero(t, sensorGot.UpdatedAt)

			sensor = sensorGot
		}

		want := entities.Tree{
			TreeCluster:         tc,
			Age:                 tree.Age,
			HeightAboveSeaLevel: tree.HeightAboveSeaLevel,
			PlantingYear:        tree.PlantingYear,
			Species:             tree.Species,
			Latitude:            tree.Latitude,
			Longitude:           tree.Longitude,
			Sensor:              sensor,
			Images:              tree.Images,
		}

		var sensorID *int32
		if sensor != nil {
			sensorID = &sensor.ID
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

		want.ID = treeGot.ID
		want.CreatedAt = treeGot.CreatedAt
		want.UpdatedAt = time.Now()

		assertTree(t, &want, treeGot)
		tc.Trees[i] = treeGot
	}
}

func TestCreateTree(t *testing.T) {
	t.Parallel()
	t.Run("should create a tree", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			// Create a tree cluster
			tc := createTreeCluster(t, str)

			// Create a tree from the tree cluster
			createTrees(t, str, tc)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMappers()
			repo := NewTreeRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.Create(context.Background(), &entities.CreateTree{
				TreeClusterID: 1,
				Age:           1,
			})
			assert.Error(t, err)
		})
	})
}

func TestGetAllTrees(t *testing.T) {
	t.Parallel()
	t.Run("should get all trees", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			tc := createTreeCluster(t, str)
			createTrees(t, str, tc)

			mappers := initMappers()
			repo := NewTreeRepository(str, mappers)

			got, err := repo.GetAll(context.Background())
			assert.NoError(t, err)
			assert.NotNil(t, got)

			assert.NotEmpty(t, got)
			assert.Len(t, got, len(tc.Trees))
			for i, tree := range tc.Trees {
				assertTree(t, tree, got[i])
			}
		})
	})

	t.Run("should return empty array if no trees found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMappers()
			repo := NewTreeRepository(str, mappers)

			got, err := repo.GetAll(context.Background())
			assert.NoError(t, err)
			assert.Empty(t, got)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMappers()
			repo := NewTreeRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.GetAll(context.Background())
			assert.Error(t, err)
		})
	})
}

func TestGetByID(t *testing.T) {
	t.Parallel()
	t.Run("should get tree by id", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			tc := createTreeCluster(t, str)
			createTrees(t, str, tc)

			mappers := initMappers()
			repo := NewTreeRepository(str, mappers)

			for _, tree := range tc.Trees {
				got, err := repo.GetByID(context.Background(), tree.ID)
				assert.NoError(t, err)
				assert.NotNil(t, got)

				assertTree(t, tree, got)
			}
		})
	})

	t.Run("should return error if tree not found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMappers()
			repo := NewTreeRepository(str, mappers)

			got, err := repo.GetByID(context.Background(), 999)
			assert.Error(t, err)
			assert.Nil(t, got)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMappers()
			repo := NewTreeRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.GetByID(context.Background(), 1)
			assert.Error(t, err)
		})
	})
}

func TestGetByTreeClusterID(t *testing.T) {
	t.Parallel()
	t.Run("should get trees by tree cluster id", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			tc := createTreeCluster(t, str)
			createTrees(t, str, tc)

			mappers := initMappers()
			repo := NewTreeRepository(str, mappers)

			got, err := repo.GetByTreeClusterID(context.Background(), tc.ID)
			assert.NoError(t, err)
			assert.NotNil(t, got)

			assert.NotEmpty(t, got)
			assert.Len(t, got, len(tc.Trees))
			for i, tree := range tc.Trees {
				assertTree(t, tree, got[i])
			}
		})
	})

	t.Run("should return empty array if no trees found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMappers()
			repo := NewTreeRepository(str, mappers)

			got, err := repo.GetByTreeClusterID(context.Background(), 999)
			assert.NoError(t, err)
			assert.Empty(t, got)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMappers()
			repo := NewTreeRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.GetByTreeClusterID(context.Background(), 1)
			assert.Error(t, err)
		})
	})
}

func TestUpdateTree(t *testing.T) {
	t.Parallel()
	t.Run("should update tree", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			tc := createTreeCluster(t, str)
			createTrees(t, str, tc)

			mappers := initMappers()
			repo := NewTreeRepository(str, mappers)

			for _, tree := range tc.Trees {
				tree.Age = 10
				tree.HeightAboveSeaLevel = 10.1
				tree.PlantingYear = 2025
				tree.Species = "updated species"
				tree.Latitude = 10.2
				tree.Longitude = 10.3

				arg := &entities.UpdateTree{
					ID:                  tree.ID,
					Age:                 utils.P(tree.Age),
					HeightAboveSeaLevel: utils.P(tree.HeightAboveSeaLevel),
					PlantingYear:        utils.P(tree.PlantingYear),
					Species:             utils.P(tree.Species),
					Latitude:            utils.P(tree.Latitude),
					Longitude:           utils.P(tree.Longitude),
				}

				got, err := repo.Update(context.Background(), arg)
				assert.NoError(t, err)
				assert.NotNil(t, got)

				assertTree(t, tree, got)
			}
		})
	})

	t.Run("should only update present fields", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			tc := createTreeCluster(t, str)
			createTrees(t, str, tc)

			mappers := initMappers()
			repo := NewTreeRepository(str, mappers)

			for _, tree := range tc.Trees {
				tree.Age = 10
				tree.HeightAboveSeaLevel = 10.1

				arg := &entities.UpdateTree{
					ID:                  tree.ID,
					Age:                 utils.P(tree.Age),
					HeightAboveSeaLevel: utils.P(tree.HeightAboveSeaLevel),
				}

				got, err := repo.Update(context.Background(), arg)
				assert.NoError(t, err)
				assert.NotNil(t, got)

				assertTree(t, tree, got)
			}
		})
	})

	t.Run("should update tree cluster id", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			tc := createTreeCluster(t, str)
			createTrees(t, str, tc)

			tcNew := createTreeCluster(t, str)

			mappers := initMappers()
			repo := NewTreeRepository(str, mappers)

			for _, tree := range tc.Trees {
				tree.TreeCluster = tcNew
				arg := &entities.UpdateTree{
					ID:            tree.ID,
					TreeClusterID: utils.P(tcNew.ID),
				}

				got, err := repo.Update(context.Background(), arg)
				assert.NoError(t, err)
				assert.NotNil(t, got)

				assertTree(t, tree, got)
			}
		})
	})

	t.Run("should return error if tree not found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMappers()
			repo := NewTreeRepository(str, mappers)

			arg := &entities.UpdateTree{
				ID:  999,
				Age: utils.P(int32(10)),
			}

			got, err := repo.Update(context.Background(), arg)
			assert.Error(t, err)
			assert.Nil(t, got)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMappers()
			repo := NewTreeRepository(str, mappers)

			arg := &entities.UpdateTree{
				ID:  1,
				Age: utils.P(int32(10)),
			}

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.Update(context.Background(), arg)
			assert.Error(t, err)
		})
	})
}

func TestDeleteTree(t *testing.T) {
	t.Parallel()
	t.Run("should delete tree", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			tc := createTreeCluster(t, str)
			createTrees(t, str, tc)

			mappers := initMappers()
			repo := NewTreeRepository(str, mappers)

			for _, tree := range tc.Trees {
				err := repo.Delete(context.Background(), tree.ID)
				assert.NoError(t, err)

				got, err := repo.GetByID(context.Background(), tree.ID)
				assert.Error(t, err)
				assert.Nil(t, got)
			}
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMappers()
			repo := NewTreeRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			err = repo.Delete(context.Background(), 1)
			assert.Error(t, err)
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
