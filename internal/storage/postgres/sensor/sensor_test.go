package sensor

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/testutils"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func createStore(db *pgx.Conn) *store.Store {
	return store.NewStore(db)
}

type RandomSensor struct {
	ID        int32                 `faker:"-"`
	CreatedAt time.Time             `faker:"-"`
	UpdatedAt time.Time             `faker:"-"`
	Status    entities.SensorStatus `faker:"oneof:online,offline,unknown"`
}

func TestMain(m *testing.M) {
	closeCon, _, err := testutils.SetupPostgresContainer()
	if err != nil {
		slog.Error("Error setting up postgres container", "error", err)
		os.Exit(1)
	}
	defer closeCon()

	os.Exit(m.Run())
}

func createSensor(t *testing.T, str *store.Store) *entities.Sensor {
	var sensor RandomSensor
	if err := faker.FakeData(&sensor); err != nil {
		t.Fatalf("error faking sensor data: %v", err)
	}
	mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
	repo := NewSensorRepository(str, mappers)

	got, err := repo.Create(context.Background(), &entities.CreateSensor{
		Status: sensor.Status,
	})
	assert.NoError(t, err)
	assert.NotNil(t, got)

	want := &entities.Sensor{
		ID:        got.ID,
		Status:    sensor.Status,
		CreatedAt: got.CreatedAt,
		UpdatedAt: got.UpdatedAt,
	}

	assertSensor(t, got, want)
	return got
}

func TestCreateSensor(t *testing.T) {
	t.Parallel()
	t.Run("should create sensor", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			createSensor(t, str)
		})
	})

	t.Run("should return error if status is invalid", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			_, err := repo.Create(context.Background(), &entities.CreateSensor{
				Status: "invalid",
			})
			assert.Error(t, err)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.Create(context.Background(), &entities.CreateSensor{
				Status: entities.SensorStatusOnline,
			})
			assert.Error(t, err)
		})
	})
}

func TestGetAllSensors(t *testing.T) {
	t.Parallel()
	t.Run("should get all sensors", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			createSensor(t, str)
			createSensor(t, str)
			createSensor(t, str)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			got, err := repo.GetAll(context.Background())
			assert.NoError(t, err)
			assert.NotEmpty(t, got)
			assert.Len(t, got, 3)
			for _, sensor := range got {
				assertSensor(t, sensor, sensor)
			}
		})
	})

	t.Run("should return empty list if no sensors", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			got, err := repo.GetAll(context.Background())
			assert.NoError(t, err)
			assert.Empty(t, got)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.GetAll(context.Background())
			assert.Error(t, err)
		})
	})
}

func TestGetSensorByID(t *testing.T) {
	t.Parallel()
	t.Run("should get sensor by id", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			sensor := createSensor(t, str)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			got, err := repo.GetByID(context.Background(), sensor.ID)
			assert.NoError(t, err)
			assertSensor(t, got, sensor)
		})
	})

	t.Run("should return error if sensor not found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			_, err := repo.GetByID(context.Background(), 999)
			assert.Error(t, err)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.GetByID(context.Background(), 1)
			assert.Error(t, err)
		})
	})
}

func TestGetStatusByID(t *testing.T) {
	t.Parallel()
	t.Run("should get sensor status by id", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			sensor := createSensor(t, str)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			got, err := repo.GetStatusByID(context.Background(), sensor.ID)
			assert.NoError(t, err)
			assert.Equal(t, got, &sensor.Status)
		})
	})

	t.Run("should return error if sensor not found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			_, err := repo.GetStatusByID(context.Background(), 999)
			assert.Error(t, err)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.GetStatusByID(context.Background(), 1)
			assert.Error(t, err)
		})
	})
}

func TestGetSensorByStatus(t *testing.T) {
	t.Parallel()
	t.Run("should get sensors by status", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			createdSensors := []*entities.Sensor{
				createSensor(t, str),
				createSensor(t, str),
				createSensor(t, str),
			}

			var onlineSensors []*entities.Sensor
			var offlineSensors []*entities.Sensor
			var unknownSensors []*entities.Sensor
			for _, sensor := range createdSensors {
				switch sensor.Status {
				case entities.SensorStatusOnline:
					onlineSensors = append(onlineSensors, sensor)
				case entities.SensorStatusOffline:
					offlineSensors = append(offlineSensors, sensor)
				case entities.SensorStatusUnknown:
					unknownSensors = append(unknownSensors, sensor)
				}
			}

			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			gotOnline, err := repo.GetSensorByStatus(context.Background(), utils.P(entities.SensorStatusOnline))
			assert.NoError(t, err)

			gotOffline, err := repo.GetSensorByStatus(context.Background(), utils.P(entities.SensorStatusOffline))
			assert.NoError(t, err)

			gotUnknown, err := repo.GetSensorByStatus(context.Background(), utils.P(entities.SensorStatusUnknown))
			assert.NoError(t, err)

			if len(onlineSensors) > 0 {
				assert.Len(t, gotOnline, len(onlineSensors))
				for i, sensor := range gotOnline {
					assertSensor(t, sensor, onlineSensors[i])
				}
			} else {
				assert.Empty(t, gotOnline)
			}

			if len(offlineSensors) > 0 {
				assert.Len(t, gotOffline, len(offlineSensors))
				for i, sensor := range gotOffline {
					assertSensor(t, sensor, offlineSensors[i])
				}
			} else {
				assert.Empty(t, gotOffline)
			}

			if len(unknownSensors) > 0 {
				assert.Len(t, gotUnknown, len(unknownSensors))
				for i, sensor := range gotUnknown {
					assertSensor(t, sensor, unknownSensors[i])
				}
			} else {
				assert.Empty(t, gotUnknown)
			}
		})
	})

	t.Run("should return empty list if no sensors", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			gotOnline, err := repo.GetSensorByStatus(context.Background(), utils.P(entities.SensorStatusOnline))
			assert.NoError(t, err)

			gotOffline, err := repo.GetSensorByStatus(context.Background(), utils.P(entities.SensorStatusOffline))
			assert.NoError(t, err)

			gotUnknown, err := repo.GetSensorByStatus(context.Background(), utils.P(entities.SensorStatusUnknown))
			assert.NoError(t, err)

			assert.Empty(t, gotOnline)
			assert.Empty(t, gotOffline)
			assert.Empty(t, gotUnknown)
		})
	})

	t.Run("should return error if status is invalid", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			_, err := repo.GetSensorByStatus(context.Background(), utils.P(entities.SensorStatus("invalid")))
			assert.Error(t, err)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.GetSensorByStatus(context.Background(), utils.P(entities.SensorStatusOnline))
			assert.Error(t, err)
		})
	})
}

func TestGetSensorDataByID(t *testing.T) {
	t.Skip("should get sensor data by id")

	t.Run("should return error if sensor not found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			_, err := repo.GetSensorDataByID(context.Background(), 999)
			assert.Error(t, err)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.GetSensorDataByID(context.Background(), 1)
			assert.Error(t, err)
		})
	})
}

func TestInsertSensorData(t *testing.T) {
	t.Skip("should insert sensor data")

	t.Run("should return error if sensor not found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			data := make([]*entities.SensorData, 1)
			data[0] = &entities.SensorData{
				ID: 999,
			}
			_, err := repo.InsertSensorData(context.Background(), data)
			assert.Error(t, err)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			data := make([]*entities.SensorData, 1)
			data[0] = &entities.SensorData{
				ID: 1,
			}
			_, err = repo.InsertSensorData(context.Background(), data)
			assert.Error(t, err)
		})
	})
}

func TestUpdateSensor(t *testing.T) {
	t.Run("should update sensor", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)

			prev := createSensor(t, str)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)
			want := &entities.Sensor{
				ID:        prev.ID,
				Status:    entities.SensorStatusOffline,
				CreatedAt: prev.CreatedAt,
				UpdatedAt: time.Now(),
			}

			got, err := repo.Update(context.Background(), &entities.UpdateSensor{
				ID:     prev.ID,
				Status: entities.SensorStatusOffline,
			})

			assert.NoError(t, err)
			assertSensor(t, got, want)
		})
	})

	t.Run("should return error if sensor not found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			_, err := repo.Update(context.Background(), &entities.UpdateSensor{
				ID:     999,
				Status: entities.SensorStatusOffline,
			})
			assert.Error(t, err)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.Update(context.Background(), &entities.UpdateSensor{
				ID:     1,
				Status: entities.SensorStatusOffline,
			})
			assert.Error(t, err)
		})
	})
}

func TestDeleteSensor(t *testing.T) {
	t.Run("should delete sensor", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			prev := createSensor(t, str)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			err := repo.Delete(context.Background(), prev.ID)
			assert.NoError(t, err)

			_, err = repo.GetByID(context.Background(), prev.ID)
			assert.Error(t, err)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := NewSensorRepositoryMappers(&generated.InternalSensorRepoMapperImpl{})
			repo := NewSensorRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			err = repo.Delete(context.Background(), 1)
			assert.Error(t, err)
		})
	})
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
