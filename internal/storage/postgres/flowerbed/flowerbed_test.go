package flowerbed

import (
	"context"
	"log/slog"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/testutils"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"

	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	mapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/generated"
)

type RandomFlowerbed struct {
	ID             int32          `faker:"-"`
	CreatedAt      time.Time      `faker:"-"`
	UpdatedAt      time.Time      `faker:"-"`
	Size           float64        `faker:"oneof:1.5,2.5,3.5"`
	Description    string         `faker:"paragraph"`
	NumberOfPlants int32          `faker:"oneof:1,2,3"`
	MoistureLevel  float64        `faker:"oneof:0.5,0.6,0.7"`
	Region         string         `faker:"oneof:Neustadt,Mürwik,Jürgensby"`
	Address        string         `faker:"real_address"`
	Sensor         *RandomSensor  `faker:"randomSensor"`
	Images         []*RandomImage `faker:"randomImage"`
	Archived       bool           `faker:"-"`
	Latitude       float64        `faker:"lat"`
	Longitude      float64        `faker:"long"`
}

type RandomSensor struct {
	ID        int32                 `faker:"-"`
	CreatedAt time.Time             `faker:"-"`
	UpdatedAt time.Time             `faker:"-"`
	Status    entities.SensorStatus `faker:"oneof:online,offline,unknown"`
}

type RandomImage struct {
	ID        int32     `faker:"-"`
	CreatedAt time.Time `faker:"-"`
	UpdatedAt time.Time `faker:"-"`
	URL       string    `faker:"url"`
	Filename  *string   `faker:"word"`
	MimeType  *string   `faker:"oneof:image/png,image/jpeg"`
}

func initFaker() {
	onErr := func(err error) {
		slog.Error("Error faking data", "error", err)
		os.Exit(1)
	}

	err := faker.AddProvider("randomImage", func(v reflect.Value) (interface{}, error) {
		images := make([]*RandomImage, 3)
		for i := 0; i < 3; i++ {
			img := RandomImage{}
			err := faker.FakeData(&img)
			if err != nil {
				return nil, err
			}
			images[i] = &img
		}

		return images, nil
	})

	if err != nil {
		onErr(err)
	}

	err = faker.AddProvider("randomSensor", func(v reflect.Value) (interface{}, error) {
		sensor := RandomSensor{}
		err := faker.FakeData(&sensor)
		if err != nil {
			return nil, err
		}

		return &sensor, nil
	})

	if err != nil {
		onErr(err)
	}

}

func mapperRepo() FlowerbedMappers {
	return NewFlowerbedMappers(&mapper.InternalFlowerbedRepoMapperImpl{}, &mapper.InternalImageRepoMapperImpl{}, &mapper.InternalSensorRepoMapperImpl{})
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

func createFlowerbed(t *testing.T, str *store.Store) *entities.Flowerbed {
	var want RandomFlowerbed
	err := faker.FakeData(&want)
	if err != nil {
		t.Fatal(err)
	}
	repo := NewFlowerbedRepository(str, mapperRepo())

	// Create sensor
	sensorRepo := sensor.NewSensorRepository(str, sensor.NewSensorRepositoryMappers(&mapper.InternalSensorRepoMapperImpl{}))
	sensorGot, err := sensorRepo.Create(context.Background(), sensor.WithStatus(want.Sensor.Status))
	assert.NoError(t, err)
	want.Sensor.ID = sensorGot.ID

	// Create images
	for i, img := range want.Images {
		args := sqlc.CreateImageParams{
			Url:      img.URL,
			Filename: img.Filename,
			MimeType: img.MimeType,
		}
		imgID, err := str.CreateImage(context.Background(), &args)
		want.Images[i].ID = imgID

		assert.NoError(t, err)
	}

	got, err := repo.Create(context.Background(),
		WithSize(want.Size),
		WithDescription(want.Description),
		WithNumberOfPlants(want.NumberOfPlants),
		WithMoistureLevel(want.MoistureLevel),
		WithRegion(want.Region),
		WithAddress(want.Address),
		WithArchived(want.Archived),
		WithLatitude(want.Latitude),
		WithLongitude(want.Longitude),
		WithSensor(&entities.Sensor{ID: want.Sensor.ID}),
		WithImages(want.Images),
	)
	assert.NoError(t, err)
	want.ID = got.ID

	assert.NotZero(t, got.CreatedAt)
	assert.NotZero(t, got.UpdatedAt)

	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.Size, got.Size)
	assert.Equal(t, want.Description, got.Description)
	assert.Equal(t, want.NumberOfPlants, got.NumberOfPlants)
	assert.Equal(t, want.MoistureLevel, got.MoistureLevel)
	assert.Equal(t, want.Region, got.Region)
	assert.Equal(t, want.Address, got.Address)
	assert.Equal(t, want.Archived, got.Archived)
	assert.Equal(t, want.Latitude, got.Latitude)
	assert.Equal(t, want.Longitude, got.Longitude)

	assert.NotZero(t, got.Sensor.CreatedAt)
	assert.NotZero(t, got.Sensor.UpdatedAt)

	assert.Equal(t, got.Sensor.ID, want.Sensor.ID)
	assert.Equal(t, got.Sensor.Status, want.Sensor.Status)

	assert.Len(t, got.Images, len(want.Images))
	for i := range got.Images {
		gImg := got.Images[i]
		wImg := want.Images[i]
		assert.Equal(t, wImg.ID, gImg.ID)
		assert.Equal(t, wImg.URL, gImg.URL)
		assert.Equal(t, wImg.Filename, gImg.Filename)
		assert.Equal(t, wImg.MimeType, gImg.MimeType)
	}

	assert.NotZero(t, got.Sensor.CreatedAt)
	assert.NotZero(t, got.Sensor.UpdatedAt)

	assert.Equal(t, want.Sensor.ID, got.Sensor.ID)
	assert.Equal(t, want.Sensor.Status, got.Sensor.Status)

	return got
}

func TestCreateFlowerbed(t *testing.T) {
	t.Parallel()
	t.Run("should create flowerbed", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			createFlowerbed(t, str)
		})
	})

	t.Run("should return error when sensor not found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			repo := NewFlowerbedRepository(str, mapperRepo())

			sensorID := int32(999)
			_, err := repo.Create(context.Background(), &entities.CreateFlowerbed{
				SensorID: &sensorID,
			})
			assert.Error(t, err)

			// FIXME: Check if entity is not created
		})
	})

	t.Run("should return error when image not found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			repo := NewFlowerbedRepository(str, mapperRepo())
			imageID := int32(999)
			_, err := repo.Create(context.Background(), &entities.CreateFlowerbed{
				ImageIDs: []int32{imageID},
			})
			assert.Error(t, err)

			// FIXME: Check if entity is not created

		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			repo := NewFlowerbedRepository(str, mapperRepo())

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.Create(context.Background(), &entities.CreateFlowerbed{
				Size: 1.5,
			})
			assert.Error(t, err)
		})
	})
}

func TestGetAll(t *testing.T) {
	t.Parallel()
	t.Run("should get all flowerbeds", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			w1 := createFlowerbed(t, str)
			w2 := createFlowerbed(t, str)
			w3 := createFlowerbed(t, str)

			repo := NewFlowerbedRepository(str, mapperRepo())
			got, err := repo.GetAll(context.Background())
			assert.NoError(t, err)

			assert.Len(t, got, 3)
			assert.Contains(t, got, w1)
			assert.Contains(t, got, w2)
			assert.Contains(t, got, w3)
		})
	})

	t.Run("should return empty list", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			repo := NewFlowerbedRepository(str, mapperRepo())

			got, err := repo.GetAll(context.Background())

			assert.NoError(t, err)
			assert.Empty(t, got)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			repo := NewFlowerbedRepository(str, mapperRepo())

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.GetAll(context.Background())
			assert.Error(t, err)
		})
	})
}

func TestGetByID(t *testing.T) {
	t.Parallel()
	t.Run("should get flowerbed by id", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			want := createFlowerbed(t, str)

			repo := NewFlowerbedRepository(str, mapperRepo())
			got, err := repo.GetByID(context.Background(), want.ID)

			assert.NoError(t, err)
			assert.Equal(t, want, got)
		})
	})

	t.Run("should return error when flowerbed not found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			repo := NewFlowerbedRepository(str, mapperRepo())

			_, err := repo.GetByID(context.Background(), 999)
			assert.Error(t, err)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			repo := NewFlowerbedRepository(str, mapperRepo())

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.GetByID(context.Background(), 1)
			assert.Error(t, err)
		})
	})
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	t.Run("should update flowerbed", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			prev := createFlowerbed(t, str)
			want := prev
			want.Images = []*entities.Image{
				{
					ID:        prev.Images[0].ID,
					URL:       prev.Images[0].URL,
					Filename:  prev.Images[0].Filename,
					MimeType:  prev.Images[0].MimeType,
					CreatedAt: prev.Images[0].CreatedAt,
					UpdatedAt: prev.Images[0].UpdatedAt,
				},
				{
					ID:        prev.Images[1].ID,
					URL:       prev.Images[1].URL,
					Filename:  prev.Images[1].Filename,
					MimeType:  prev.Images[1].MimeType,
					CreatedAt: prev.Images[1].CreatedAt,
					UpdatedAt: prev.Images[1].UpdatedAt,
				},
			}

			repo := NewFlowerbedRepository(str, mapperRepo())
			got, err := repo.Update(context.Background(), &entities.UpdateFlowerbed{
				ID:             want.ID,
				Size:           &want.Size,
				Description:    &want.Description,
				NumberOfPlants: &want.NumberOfPlants,
				MoistureLevel:  &want.MoistureLevel,
				Region:         &want.Region,
				Address:        &want.Address,
				Archived:       &want.Archived,
				Latitude:       &want.Latitude,
				Longitude:      &want.Longitude,
				SensorID:       &want.Sensor.ID,
				ImageIDs:       []int32{want.Images[0].ID, want.Images[1].ID}, // should remove the third image
			})

			assert.NoError(t, err)
			assert.WithinDuration(t, time.Now(), got.UpdatedAt, time.Second)
			assertFlowerbed(t, got, want)
		})
	})

	t.Run("should only update fields that are not nil", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			prev := createFlowerbed(t, str)

			want := RandomFlowerbed{}
			err := faker.FakeData(&want)
			if err != nil {
				t.Fatal(err)
			}
			want.ID = prev.ID

			repo := NewFlowerbedRepository(str, mapperRepo())
			got, err := repo.Update(context.Background(), &entities.UpdateFlowerbed{
				ID:             want.ID,
				Size:           &want.Size,           // should update size
				Description:    &want.Description,    // should update description
				NumberOfPlants: &want.NumberOfPlants, // should update number of plants
				MoistureLevel:  nil,                  // should not update moisture level
				Region:         nil,                  // should not update region
				Address:        nil,                  // should not update address
				Archived:       nil,                  // should not update archived
				Latitude:       nil,                  // should not update location
				Longitude:      nil,                  // should not update location
				SensorID:       nil,                  // should not update sensor
				ImageIDs:       nil,                  // should not update images
			})

			assert.NoError(t, err)
			assert.WithinDuration(t, time.Now(), got.UpdatedAt, time.Second)
			assert.Equal(t, prev.ID, got.ID)
			assert.Equal(t, want.Size, got.Size)
			assert.Equal(t, want.Description, got.Description)
			assert.Equal(t, want.NumberOfPlants, got.NumberOfPlants)
			assert.Equal(t, prev.MoistureLevel, got.MoistureLevel)
			assert.Equal(t, prev.Region, got.Region)
			assert.Equal(t, prev.Address, got.Address)
			assert.Equal(t, prev.Archived, got.Archived)
			assert.Equal(t, prev.Latitude, got.Latitude)
			assert.Equal(t, prev.Longitude, got.Longitude)

			assertSensor(t, prev.Sensor, got.Sensor)

			assert.Len(t, got.Images, len(prev.Images))
			assertImage(t, prev.Images[0], got.Images[0])
			assertImage(t, prev.Images[1], got.Images[1])
			assertImage(t, prev.Images[2], got.Images[2])
		})
	})

	t.Run("should return error when flowerbed not found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			repo := NewFlowerbedRepository(str, mapperRepo())

			_, err := repo.Update(context.Background(), &entities.UpdateFlowerbed{
				ID: 999,
			})
			assert.Error(t, err)
		})
	})

	t.Run("should return error when sensor not found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			repo := NewFlowerbedRepository(str, mapperRepo())

			_, err := repo.Update(context.Background(), &entities.UpdateFlowerbed{
				ID:       1,
				SensorID: new(int32),
			})
			assert.Error(t, err)
		})
	})

	t.Run("should return error when image not found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			repo := NewFlowerbedRepository(str, mapperRepo())

			_, err := repo.Update(context.Background(), &entities.UpdateFlowerbed{
				ID:       1,
				ImageIDs: []int32{999},
			})
			assert.Error(t, err)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			repo := NewFlowerbedRepository(str, mapperRepo())

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.Update(context.Background(), &entities.UpdateFlowerbed{ID: 1})
			assert.Error(t, err)
		})
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()
	t.Run("should delete flowerbed", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			want := createFlowerbed(t, str)

			repo := NewFlowerbedRepository(str, mapperRepo())
			err := repo.Delete(context.Background(), want.ID)

			assert.NoError(t, err)

			_, err = repo.GetByID(context.Background(), want.ID)
			assert.Error(t, err)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			repo := NewFlowerbedRepository(str, mapperRepo())

			err := db.Close(context.Background())
			assert.NoError(t, err)

			err = repo.Delete(context.Background(), 1)
			assert.Error(t, err)
		})
	})
}

func TestGetBySensorID(t *testing.T) {
	t.Parallel()
	t.Run("should get sensor by flowerbed id", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			want := createFlowerbed(t, str)

			repo := NewFlowerbedRepository(str, mapperRepo())
			got, err := repo.GetSensorByFlowerbedID(context.Background(), want.ID)

			assert.NoError(t, err)
			assertSensor(t, got, want.Sensor)
		})
	})

	t.Run("should return error when flowerbed not found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			repo := NewFlowerbedRepository(str, mapperRepo())

			_, err := repo.GetSensorByFlowerbedID(context.Background(), 999)
			assert.Error(t, err)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			repo := NewFlowerbedRepository(str, mapperRepo())

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.GetSensorByFlowerbedID(context.Background(), 1)
			assert.Error(t, err)
		})
	})
}

func TestArchive(t *testing.T) {
	t.Parallel()
	t.Run("should archive flowerbed", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			want := createFlowerbed(t, str)

			repo := NewFlowerbedRepository(str, mapperRepo())
			err := repo.Archive(context.Background(), want.ID)

			assert.NoError(t, err)

			got, err := repo.GetByID(context.Background(), want.ID)
			assert.NoError(t, err)
			assert.True(t, got.Archived)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := store.NewStore(db)
			repo := NewFlowerbedRepository(str, mapperRepo())

			err := db.Close(context.Background())
			assert.NoError(t, err)

			err = repo.Archive(context.Background(), 1)
			assert.Error(t, err)
		})
	})
}

func assertFlowerbed(t *testing.T, got, want *entities.Flowerbed) {
	if got == nil {
		assert.Nil(t, got)
		return
	}

	if want == nil {
		assert.Nil(t, want)
		return
	}

	assert.NotZero(t, got.CreatedAt)
	assert.NotZero(t, got.UpdatedAt)

	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.Size, got.Size)
	assert.Equal(t, want.Description, got.Description)
	assert.Equal(t, want.NumberOfPlants, got.NumberOfPlants)
	assert.Equal(t, want.MoistureLevel, got.MoistureLevel)
	assert.Equal(t, want.Region, got.Region)
	assert.Equal(t, want.Address, got.Address)
	assert.Equal(t, want.Archived, got.Archived)
	assert.Equal(t, want.Latitude, got.Latitude)
	assert.Equal(t, want.Longitude, got.Longitude)

	assert.NotZero(t, got.Sensor.CreatedAt)
	assert.NotZero(t, got.Sensor.UpdatedAt)
	assert.Equal(t, got.Sensor.ID, want.Sensor.ID)
	assert.Equal(t, got.Sensor.Status, want.Sensor.Status)

	assert.Len(t, got.Images, len(want.Images))
	for i := range got.Images {
		assertImage(t, got.Images[i], want.Images[i])
	}

	assertSensor(t, got.Sensor, want.Sensor)
}

func assertSensor(t *testing.T, got, want *entities.Sensor) {
	if got == nil {
		assert.Nil(t, got)
		return
	}

	if want == nil {
		assert.Nil(t, want)
		return
	}

	assert.NotZero(t, got.CreatedAt)
	assert.NotZero(t, got.UpdatedAt)

	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.Status, got.Status)
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
