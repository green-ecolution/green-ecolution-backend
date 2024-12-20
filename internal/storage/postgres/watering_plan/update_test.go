package wateringplan

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestWateringPlanRepository_Update(t *testing.T) {
	suite.ResetDB(t)
	suite.InsertSeed(t, "internal/storage/postgres/seed/test/watering_plan")

	testVehicles, err := suite.Store.GetAllVehicles(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	testCluster, err := suite.Store.GetAllTreeClusters(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	uuid, _ := uuid.NewRandom()
	testUser := &entities.User{
		ID:          uuid,
		CreatedAt:   time.Unix(123456, 0),
		Username:    "test",
		FirstName:   "Toni",
		LastName:    "Tester",
		Email:       "dev@green-ecolution.de",
		PhoneNumber: "+49 123456",
		EmployeeID:  "123456",
	}

	input := entities.WateringPlan{
		Date:               time.Date(2024, 11, 22, 0, 0, 0, 0, time.UTC),
		Description:        "Updated watering plan",
		Distance:           utils.P(50.0),
		TotalWaterRequired: utils.P(30000.0),
		Trailer:            mappers.vehicleMapper.FromSqlList(testVehicles)[2],
		Transporter:        mappers.vehicleMapper.FromSqlList(testVehicles)[3],
		TreeClusters:       mappers.clusterMapper.FromSqlList(testCluster)[0:3],
		Users:              []*entities.User{testUser},
		Status:             entities.WateringPlanStatusActive,
	}

	t.Run("should update watering plan", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		updateFn := func(wp *entities.WateringPlan) (bool, error) {
			wp.Date = input.Date
			wp.Description = input.Description
			wp.Distance = input.Distance
			wp.TotalWaterRequired = input.TotalWaterRequired
			wp.Transporter = input.Transporter
			wp.Trailer = input.Trailer
			wp.TreeClusters = input.TreeClusters
			wp.Users = input.Users
			wp.Status = input.Status
			return true, nil
		}

		// when
		updateErr := r.Update(context.Background(), 1, updateFn)
		got, getErr := r.GetByID(context.Background(), 1)

		// then
		assert.NoError(t, updateErr)
		assert.NoError(t, getErr)
		assert.NotNil(t, got)
		assert.NotZero(t, got.ID)
		assert.Equal(t, input.Date, got.Date)
		assert.Equal(t, input.Description, got.Description)
		assert.Equal(t, input.Distance, got.Distance)
		assert.Equal(t, input.TotalWaterRequired, got.TotalWaterRequired)
		assert.Equal(t, input.Status, got.Status)

		// assert transporter
		assert.Equal(t, input.Transporter.ID, got.Transporter.ID)
		assert.Equal(t, input.Transporter.NumberPlate, got.Transporter.NumberPlate)

		// assert trailer
		assert.Equal(t, input.Trailer.ID, got.Trailer.ID)
		assert.Equal(t, input.Trailer.NumberPlate, got.Trailer.NumberPlate)

		// assert TreeClusters
		assert.Len(t, input.TreeClusters, len(got.TreeClusters))
		for i, tc := range got.TreeClusters {
			assert.Equal(t, input.TreeClusters[i].ID, tc.ID)
			assert.Equal(t, input.TreeClusters[i].Name, tc.Name)
		}

		// TODO: test linked users
	})

	t.Run("should update watering plan and unlink trailer", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		updateFn := func(wp *entities.WateringPlan) (bool, error) {
			wp.Date = input.Date
			wp.Description = input.Description
			wp.Distance = input.Distance
			wp.TotalWaterRequired = input.TotalWaterRequired
			wp.Transporter = input.Transporter
			wp.Trailer = nil
			wp.TreeClusters = input.TreeClusters
			wp.Users = input.Users
			wp.Status = input.Status
			return true, nil
		}

		// when
		updateErr := r.Update(context.Background(), 2, updateFn)
		got, getErr := r.GetByID(context.Background(), 2)

		// then
		assert.NoError(t, updateErr)
		assert.NoError(t, getErr)
		assert.NotNil(t, got)
		assert.NotZero(t, got.ID)
		assert.Equal(t, input.Date, got.Date)
		assert.Equal(t, input.Description, got.Description)
		assert.Equal(t, input.Distance, got.Distance)
		assert.Equal(t, input.TotalWaterRequired, got.TotalWaterRequired)
		assert.Equal(t, input.Status, got.Status)

		// assert transporter
		assert.Equal(t, input.Transporter.ID, got.Transporter.ID)
		assert.Equal(t, input.Transporter.NumberPlate, got.Transporter.NumberPlate)

		// assert nil trailer
		assert.Nil(t, got.Trailer)

		// assert TreeClusters
		assert.Len(t, input.TreeClusters, len(got.TreeClusters))
		for i, tc := range got.TreeClusters {
			assert.Equal(t, input.TreeClusters[i].ID, tc.ID)
			assert.Equal(t, input.TreeClusters[i].Name, tc.Name)
		}

		// TODO: test linked users
	})

	t.Run("should update watering plan to canceled", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		canellationNote := "This watering plan is canceled"

		updateFn := func(wp *entities.WateringPlan) (bool, error) {
			wp.Date = input.Date
			wp.Description = input.Description
			wp.Distance = input.Distance
			wp.TotalWaterRequired = input.TotalWaterRequired
			wp.Transporter = input.Transporter
			wp.TreeClusters = input.TreeClusters
			wp.Users = input.Users
			wp.Status = entities.WateringPlanStatusCanceled
			wp.CancellationNote = canellationNote
			return true, nil
		}

		// when
		updateErr := r.Update(context.Background(), 2, updateFn)
		got, getErr := r.GetByID(context.Background(), 2)

		// then
		assert.NoError(t, updateErr)
		assert.NoError(t, getErr)
		assert.NotNil(t, got)
		assert.NotZero(t, got.ID)
		assert.Equal(t, input.Date, got.Date)
		assert.Equal(t, entities.WateringPlanStatusCanceled, got.Status)
		assert.Equal(t, canellationNote, got.CancellationNote)
	})

	t.Run("should return error when canellation note is not empty and the status is not canceled", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		updateFn := func(wp *entities.WateringPlan) (bool, error) {
			wp.Date = input.Date
			wp.Distance = input.Distance
			wp.TotalWaterRequired = input.TotalWaterRequired
			wp.Transporter = input.Transporter
			wp.TreeClusters = input.TreeClusters
			wp.Users = input.Users
			wp.Status = entities.WateringPlanStatusActive
			wp.CancellationNote = "This watering plan is canceled"
			return true, nil
		}

		// when
		err := r.Update(context.Background(), 2, updateFn)

		// then
		assert.Error(t, err)
		assert.Equal(t, "cancellation note should be empty, as the current watering plan is not canceled", err.Error())
	})

	t.Run("should return error when date is not in correct format", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		updateFn := func(wp *entities.WateringPlan) (bool, error) {
			wp.Date = time.Time{}
			wp.Transporter = input.Transporter
			wp.TreeClusters = input.TreeClusters
			wp.Users = input.Users
			return true, nil
		}

		// when
		err := r.Update(context.Background(), 1, updateFn)

		// then
		assert.Error(t, err)
		assert.Equal(t, "failed to convert date", err.Error())
	})

	t.Run("should return error when trailer vehicle has not correct vehilce type", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		updateFn := func(wp *entities.WateringPlan) (bool, error) {
			wp.Date = input.Date
			wp.Transporter = input.Transporter
			wp.Trailer = input.Transporter
			wp.TreeClusters = input.TreeClusters
			wp.Users = input.Users
			return true, nil
		}

		// when
		err := r.Update(context.Background(), 1, updateFn)

		// then
		assert.Error(t, err)
		assert.Equal(t, "trailer vehicle requires a vehicle of type trailer", err.Error())
	})

	// t.Run("should return error when watering plan has no linked users", func(t *testing.T) {
	// 	// given
	// 	r := NewWateringPlanRepository(suite.Store, mappers)

	// 	updateFn := func(wp *entities.WateringPlan) (bool, error) {
	// 		wp.Date = input.Date
	// 		wp.Transporter = input.Transporter
	// 		wp.Trailer = input.Trailer
	// 		wp.TreeClusters = input.TreeClusters
	// 		wp.Users = []*entities.User{}
	// 		return true, nil
	// 	}

	// 	// when
	// 	err := r.Update(context.Background(), 1, updateFn)

	// 	// then
	// 	assert.Error(t, err)
	// 	assert.Equal(t, "watering plan requires employees", err.Error())
	// })

	t.Run("should return error when transporter has not correct vehilce type", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		updateFn := func(wp *entities.WateringPlan) (bool, error) {
			wp.Date = input.Date
			wp.Transporter = input.Trailer
			wp.Trailer = input.Trailer
			wp.TreeClusters = input.TreeClusters
			wp.Users = input.Users
			return true, nil
		}

		// when
		err := r.Update(context.Background(), 1, updateFn)

		// then
		assert.Error(t, err)
		assert.Equal(t, "watering plan requires a valid transporter", err.Error())
	})

	t.Run("should return error when transporter is nil", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		updateFn := func(wp *entities.WateringPlan) (bool, error) {
			wp.Date = input.Date
			wp.Transporter = nil
			wp.Trailer = input.Trailer
			wp.TreeClusters = input.TreeClusters
			wp.Users = input.Users
			return true, nil
		}

		// when
		err := r.Update(context.Background(), 1, updateFn)

		// then
		assert.Error(t, err)
		assert.Equal(t, "watering plan requires a valid transporter", err.Error())
	})

	t.Run("should return error when no TreeClusters are linked", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		updateFn := func(wp *entities.WateringPlan) (bool, error) {
			wp.Date = input.Date
			wp.Transporter = input.Transporter
			wp.Trailer = input.Trailer
			wp.TreeClusters = []*entities.TreeCluster{}
			wp.Users = input.Users
			return true, nil
		}

		// when
		err := r.Update(context.Background(), 1, updateFn)

		// then
		assert.Error(t, err)
		assert.Equal(t, "watering plan requires tree cluster", err.Error())
	})

	t.Run("should return error when watering plan is invalid", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		err := r.Update(context.Background(), 1, nil)

		// then
		assert.Error(t, err)
		assert.Equal(t, "updateFn is nil", err.Error())
	})

	t.Run("should return error when update watering plan with negative id", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		updateFn := func(wp *entities.WateringPlan) (bool, error) {
			wp.Date = input.Date
			wp.Transporter = input.Transporter
			wp.Trailer = input.Trailer
			wp.TreeClusters = input.TreeClusters
			wp.Users = input.Users
			return true, nil
		}

		// when
		err := r.Update(context.Background(), -1, updateFn)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when update watering plan with zero id", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		updateFn := func(wp *entities.WateringPlan) (bool, error) {
			wp.Date = input.Date
			wp.Transporter = input.Transporter
			wp.Trailer = input.Trailer
			wp.TreeClusters = input.TreeClusters
			wp.Users = input.Users
			return true, nil
		}

		// when
		err := r.Update(context.Background(), 0, updateFn)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when update watering plan not found", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		updateFn := func(wp *entities.WateringPlan) (bool, error) {
			wp.Date = input.Date
			wp.Transporter = input.Transporter
			wp.Trailer = input.Trailer
			wp.TreeClusters = input.TreeClusters
			wp.Users = input.Users
			return true, nil
		}

		// when
		err := r.Update(context.Background(), 99, updateFn)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error if context is canceled", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)
		updateFn := func(wp *entities.WateringPlan) (bool, error) {
			wp.Date = input.Date
			wp.Transporter = input.Transporter
			wp.Trailer = input.Trailer
			wp.TreeClusters = input.TreeClusters
			wp.Users = input.Users
			return true, nil
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		err := r.Update(ctx, 99, updateFn)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when updateFn is nil", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		err := r.Update(context.Background(), 1, nil)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when updateFn returns error", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)
		updateFn := func(wp *entities.WateringPlan) (bool, error) {
			return true, assert.AnError
		}

		// when
		err := r.Update(context.Background(), 1, updateFn)

		// then
		assert.Error(t, err)
	})

	t.Run("should not update when updateFn returns false", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)
		updateFn := func(wp *entities.WateringPlan) (bool, error) {
			return false, nil
		}

		// when
		updateErr := r.Update(context.Background(), 1, updateFn)
		got, getErr := r.GetByID(context.Background(), 1)

		// then
		assert.NoError(t, updateErr)
		assert.NoError(t, getErr)
		assert.NotNil(t, got)
	})

	t.Run("should not rollback when updateFn returns false", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)
		updateFn := func(wp *entities.WateringPlan) (bool, error) {
			wp.Date = input.Date
			wp.Description = "Test"
			wp.Transporter = input.Transporter
			wp.Trailer = input.Trailer
			wp.TreeClusters = input.TreeClusters
			wp.Users = input.Users
			return false, nil
		}

		// when
		err := r.Update(context.Background(), 1, updateFn)
		got, getErr := r.GetByID(context.Background(), 1)

		// then
		assert.NoError(t, err)
		assert.NoError(t, getErr)
		assert.NotNil(t, got)
		assert.NotEqual(t, "Test", got.Description)
	})
}
