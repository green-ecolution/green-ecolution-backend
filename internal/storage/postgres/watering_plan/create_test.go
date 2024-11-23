package watering_plan

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestWateringPlanRepository_Create(t *testing.T) {
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
		Date:               time.Date(2024, 9, 22, 0, 0, 0, 0, time.UTC),
		Description:        "New watering plan",
		Distance:           utils.P(50.0),
		TotalWaterRequired: utils.P(30000.0),
		Trailer:            mappers.vehicleMapper.FromSqlList(testVehicles)[0],
		Transporter:        mappers.vehicleMapper.FromSqlList(testVehicles)[1],
		Treecluster:        mappers.clusterMapper.FromSqlList(testCluster)[0:5],
		Users:              []*entities.User{testUser},
	}

	t.Run("should create watering plan with all values", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		got, err := r.Create(context.Background(),
			WithDate(input.Date),
			WithDescription(input.Description),
			WithDistance(input.Distance),
			WithTotalWaterRequired(input.TotalWaterRequired),
			WithTransporter(input.Transporter),
			WithTrailer(input.Trailer),
			WithTreecluster(input.Treecluster),
			WithUsers(input.Users),
		)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.NotZero(t, got.ID)
		assert.Equal(t, input.Date, got.Date)
		assert.Equal(t, input.Description, got.Description)
		assert.Equal(t, input.Distance, got.Distance)
		assert.Equal(t, input.TotalWaterRequired, got.TotalWaterRequired)

		// TODO: test linkd treeclusters, vehicles and users
	})

	t.Run("should create watering plan with default values", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		got, err := r.Create(context.Background(),
			WithDate(input.Date),
			WithTransporter(input.Transporter),
			WithTreecluster(input.Treecluster),
			WithUsers(input.Users),
		)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.NotZero(t, got.ID)
		assert.Equal(t, input.Date, got.Date)
		assert.Equal(t, "", got.Description)
		assert.Equal(t, utils.P(float64(0)), got.Distance)
		assert.Equal(t, utils.P(float64(0)), got.TotalWaterRequired)

		// TODO: test linkd treeclusters, vehicles and users
	})

	t.Run("should return error when date is not in correct format", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		got, err := r.Create(context.Background(),
			WithDate(time.Time{}),
			WithTransporter(input.Transporter),
			WithTreecluster(input.Treecluster),
			WithUsers(input.Users),
		)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.Equal(t, "failed to convert date", err.Error())
	})

	t.Run("should return error when trailer vehicle has not correct vehilce type", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		got, err := r.Create(context.Background(),
			WithDate(input.Date),
			WithTransporter(input.Transporter),
			WithTrailer(input.Transporter),
			WithTreecluster(input.Treecluster),
			WithUsers(input.Users),
		)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.Equal(t, "trailer vehicle requires a vehicle of type trailer", err.Error())
	})

	t.Run("should return error when watering plan has no linked users", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		got, err := r.Create(context.Background(),
			WithDate(input.Date),
			WithTransporter(input.Transporter),
			WithTrailer(input.Trailer),
			WithTreecluster(input.Treecluster),
			WithUsers([]*entities.User{}),
		)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.Equal(t, "watering plan requires employees", err.Error())
	})

	t.Run("should return error when transporter has not correct vehilce type", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		got, err := r.Create(context.Background(),
			WithDate(input.Date),
			WithTransporter(input.Trailer),
			WithTreecluster(input.Treecluster),
			WithUsers(input.Users),
		)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.Equal(t, "watering plan requires a valid transporter", err.Error())
	})

	t.Run("should return error when transporter is nil", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		got, err := r.Create(context.Background(),
			WithDate(input.Date),
			WithTransporter(nil),
			WithTreecluster(input.Treecluster),
			WithUsers(input.Users),
		)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.Equal(t, "watering plan requires a valid transporter", err.Error())
	})

	t.Run("should return error when no treecluster are linked", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		got, err := r.Create(context.Background(),
			WithDate(input.Date),
			WithTransporter(input.Transporter),
			WithTreecluster([]*entities.TreeCluster{}),
			WithUsers(input.Users),
		)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.Equal(t, "watering plan requires tree cluster", err.Error())
	})

	t.Run("should return error when watering plan is invalid", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		got, err := r.Create(context.Background())

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error if context is canceled", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.Create(ctx,
			WithDate(input.Date),
			WithTransporter(nil),
			WithTreecluster(input.Treecluster),
			WithUsers(input.Users),
		)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}
