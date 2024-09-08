package vehicle

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	mapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/testutils"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

type RandomVehicle struct {
	ID            int32     `faker:"-"`
	CreatedAt     time.Time `faker:"-"`
	UpdatedAt     time.Time `faker:"-"`
	NumberPlate   string    `faker:"oneof:AB1234,CD5678,EF91011"`
	Description   string    `faker:"sentence"`
	WaterCapacity float64   `faker:"float64"`
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

func createStore(db *pgx.Conn) *store.Store {
	return store.NewStore(db)
}

func initMapper() VehicleRepositoryMappers {
	return NewVehicleRepositoryMappers(&mapper.InternalVehicleRepoMapperImpl{})
}

func createVehicle(t *testing.T, str *store.Store) *entities.Vehicle {
	var v entities.Vehicle
	if err := faker.FakeData(&v); err != nil {
		t.Fatalf("error faking vehicle data: %v", err)
	}
	mappers := initMapper()
	repo := NewVehicleRepository(str, mappers)

	got, err := repo.Create(context.Background(),
		WithNumberPlate(v.NumberPlate),
		WithDescription(v.Description),
		WithWaterCapacity(v.WaterCapacity),
	)
	if err != nil {
		t.Fatalf("error creating vehicle: %v", err)
	}

	assertVehicle(t, got, &entities.Vehicle{
		ID:            got.ID,
		CreatedAt:     got.CreatedAt,
		UpdatedAt:     got.UpdatedAt,
		NumberPlate:   v.NumberPlate,
		Description:   v.Description,
		WaterCapacity: v.WaterCapacity,
	})

	return got
}

func TestCreateVehicle(t *testing.T) {
	t.Parallel()
	t.Run("should create vehicle", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			createVehicle(t, str)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMapper()
			repo := NewVehicleRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)
			_, err = repo.Create(context.Background())
			assert.Error(t, err)
		})
	})
}

func TestGetAllVehicles(t *testing.T) {
	t.Parallel()
	t.Run("should get all vehicles", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)

			createdVehicles := make([]*entities.Vehicle, 0)
			for i := 0; i < 3; i++ {
				createdVehicles = append(createdVehicles, createVehicle(t, str))
			}

			mappers := initMapper()
			repo := NewVehicleRepository(str, mappers)
			got, err := repo.GetAll(context.Background())
			assert.NoError(t, err)
			assert.Len(t, got, 3)

			for i, v := range got {
				assertVehicle(t, v, createdVehicles[i])
			}
		})
	})

	t.Run("should return empty list if no vehicles found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMapper()
			repo := NewVehicleRepository(str, mappers)

			got, err := repo.GetAll(context.Background())
			assert.NoError(t, err)
			assert.Len(t, got, 0)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMapper()
			repo := NewVehicleRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.GetAll(context.Background())
			assert.Error(t, err)
		})
	})
}

func TestGetVehicleByID(t *testing.T) {
	t.Parallel()
	t.Run("should get vehicle by id", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			v := createVehicle(t, str)
			mappers := initMapper()
			repo := NewVehicleRepository(str, mappers)

			got, err := repo.GetByID(context.Background(), v.ID)
			assert.NoError(t, err)
			assertVehicle(t, got, v)
		})
	})

	t.Run("should return error if vehicle not found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMapper()
			repo := NewVehicleRepository(str, mappers)

			_, err := repo.GetByID(context.Background(), 999)
			assert.Error(t, err)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMapper()
			repo := NewVehicleRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.GetByID(context.Background(), 1)
			assert.Error(t, err)
		})
	})
}

func TestGetVehicleByPlate(t *testing.T) {
	t.Parallel()
	t.Run("should get vehicle by plate", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			v := createVehicle(t, str)
			mappers := initMapper()
			repo := NewVehicleRepository(str, mappers)

			got, err := repo.GetByPlate(context.Background(), v.NumberPlate)
			assert.NoError(t, err)
			assertVehicle(t, got, v)
		})
	})

	t.Run("should return error if vehicle not found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMapper()
			repo := NewVehicleRepository(str, mappers)

			_, err := repo.GetByPlate(context.Background(), "AB1234")
			assert.Error(t, err)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMapper()
			repo := NewVehicleRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.GetByPlate(context.Background(), "AB1234")
			assert.Error(t, err)
		})
	})
}

func TestUpdateVehicle(t *testing.T) {
	t.Parallel()
	t.Run("should update vehicle", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			v := createVehicle(t, str)
			mappers := initMapper()
			repo := NewVehicleRepository(str, mappers)

			v.Description = "Updated description"
			v.WaterCapacity = 200.0
			v.NumberPlate = "CD5678"

			got, err := repo.Update(context.Background(), v.ID,
				WithNumberPlate(v.NumberPlate),
				WithDescription(v.Description),
				WithWaterCapacity(v.WaterCapacity),
			)
			assert.NoError(t, err)
			assertVehicle(t, got, v)
		})
	})

	t.Run("should update vehicle with only description", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			v := createVehicle(t, str)
			mappers := initMapper()
			repo := NewVehicleRepository(str, mappers)

			v.Description = "Updated description"
			got, err := repo.Update(context.Background(), v.ID, WithDescription(v.Description))

			assert.NoError(t, err)
			assertVehicle(t, got, v)
		})
	})

	t.Run("should update vehicle with only water capacity", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			v := createVehicle(t, str)
			mappers := initMapper()
			repo := NewVehicleRepository(str, mappers)

			v.WaterCapacity = 200.0
			got, err := repo.Update(context.Background(), v.ID, WithWaterCapacity(v.WaterCapacity))

			assert.NoError(t, err)
			assertVehicle(t, got, v)
		})
	})

	t.Run("should update vehicle with only number plate", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			v := createVehicle(t, str)
			mappers := initMapper()
			repo := NewVehicleRepository(str, mappers)

			v.NumberPlate = "CD5678"
			got, err := repo.Update(context.Background(), v.ID, WithNumberPlate(v.NumberPlate))

			assert.NoError(t, err)
			assertVehicle(t, got, v)
		})
	})

	t.Run("should return error if vehicle not found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMapper()

			repo := NewVehicleRepository(str, mappers)
			_, err := repo.Update(context.Background(), 999)

			assert.Error(t, err)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			v := createVehicle(t, str)
			mappers := initMapper()
			repo := NewVehicleRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.Update(context.Background(), v.ID)
			assert.Error(t, err)
		})
	})
}

func TestDeleteVehicle(t *testing.T) {
	t.Parallel()
	t.Run("should delete vehicle", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			v := createVehicle(t, str)
			mappers := initMapper()
			repo := NewVehicleRepository(str, mappers)

			err := repo.Delete(context.Background(), v.ID)
			assert.NoError(t, err)

			_, err = repo.GetByID(context.Background(), v.ID)
			assert.Error(t, err)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			v := createVehicle(t, str)
			mappers := initMapper()
			repo := NewVehicleRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			err = repo.Delete(context.Background(), v.ID)
			assert.Error(t, err)
		})
	})
}

func assertVehicle(t *testing.T, got *entities.Vehicle, want *entities.Vehicle) {
	if got == nil {
		assert.Nil(t, got)
		return
	}

	if want == nil {
		assert.Nil(t, want)
		return
	}

	assert.NotZero(t, got.ID)
	assert.NotZero(t, got.CreatedAt)
	assert.NotZero(t, got.UpdatedAt)
	assert.WithinDuration(t, time.Now(), got.UpdatedAt, time.Second)
	assert.Equal(t, want.NumberPlate, got.NumberPlate)
	assert.Equal(t, want.Description, got.Description)
	assert.Equal(t, want.WaterCapacity, got.WaterCapacity)
}
