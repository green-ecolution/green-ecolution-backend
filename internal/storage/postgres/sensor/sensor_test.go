package sensor

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor/mapper"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/test"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

var (
	querier      sqlc.Querier
  seedPath     string
  dbUrl        string
	mapperRepo   SensorRepositoryMappers
	defaultField struct {
		querier                 sqlc.Querier
		SensorRepositoryMappers SensorRepositoryMappers
	}
)

func TestMain(m *testing.M) {
	rootDir := utils.RootDir()
	seedPath := fmt.Sprintf("%s/internal/storage/postgres/test/seed/sensor", rootDir)
	close, url, err := test.SetupPostgresContainer(seedPath)
	if err != nil {
		slog.Error("Error setting up postgres container", "error", err)
		panic(err)
	}
	defer close()

  dbUrl = *url
  db, err := pgx.Connect(context.Background(), dbUrl)
  if err != nil {
    os.Exit(1)
  }
	querier = sqlc.New(db)
	mapperRepo = SensorRepositoryMappers{
		mapper: &generated.InternalSensorRepoMapperImpl{},
	}
	defaultField = struct {
		querier                 sqlc.Querier
		SensorRepositoryMappers SensorRepositoryMappers
	}{
		querier:                 querier,
		SensorRepositoryMappers: mapperRepo,
	}
	os.Exit(m.Run())
}

func TestNewSensorRepositoryMappers(t *testing.T) {
	type args struct {
		sMapper mapper.InternalSensorRepoMapper
	}
	tests := []struct {
		name string
		args args
		want SensorRepositoryMappers
	}{
		{
			name: "Test NewSensorRepositoryMappers",
			args: args{
				sMapper: &generated.InternalSensorRepoMapperImpl{},
			},
			want: SensorRepositoryMappers{
				mapper: &generated.InternalSensorRepoMapperImpl{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSensorRepositoryMappers(tt.args.sMapper); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSensorRepositoryMappers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSensorRepository(t *testing.T) {
	type args struct {
		querier sqlc.Querier
		mappers SensorRepositoryMappers
	}
	tests := []struct {
		name string
		args args
		want storage.SensorRepository
	}{
		{
			name: "Test NewSensorRepository",
			args: args{
				querier: querier,
				mappers: mapperRepo,
			},
			want: &SensorRepository{
				querier:                 querier,
				SensorRepositoryMappers: mapperRepo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSensorRepository(tt.args.querier, tt.args.mappers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSensorRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSensorRepository_GetAll(t *testing.T) {
	type fields struct {
		querier                 sqlc.Querier
		SensorRepositoryMappers SensorRepositoryMappers
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*entities.Sensor
		wantErr bool
	}{
		{
			name: "SensorRepository.GetAll should return all sensors",
			fields: fields{
				querier:                 querier,
				SensorRepositoryMappers: mapperRepo,
			},
			args: args{
				ctx: context.Background(),
			},
			want: []*entities.Sensor{
				{
					ID:     1,
					Status: entities.SensorStatusOnline,
				},
				{
					ID:     2,
					Status: entities.SensorStatusOffline,
				},
				{
					ID:     3,
					Status: entities.SensorStatusUnknown,
				},
				{
					ID:     4,
					Status: entities.SensorStatusOnline,
				},
			},
		},
	}
	for _, tt := range tests {
    test.ResetDatabase(dbUrl, seedPath)
		t.Run(tt.name, func(t *testing.T) {
			r := &SensorRepository{
				querier:                 tt.fields.querier,
				SensorRepositoryMappers: tt.fields.SensorRepositoryMappers,
			}
			got, err := r.GetAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Len(t, got, len(tt.want))
			for i := range got {
				assertSensor(t, got[i], tt.want[i])
			}
		})
	}
}

func TestSensorRepository_GetByID(t *testing.T) {
	type fields struct {
		querier                 sqlc.Querier
		SensorRepositoryMappers SensorRepositoryMappers
	}
	type args struct {
		ctx context.Context
		id  int32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entities.Sensor
		wantErr bool
	}{
		{
			name:   "SensorRepository.GetByID should return sensor by id",
			fields: defaultField,
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: &entities.Sensor{
				ID:     1,
				Status: entities.SensorStatusOnline,
			},
		},
		{
			name:   "SensorRepository.GetByID should return error when id not found",
			fields: defaultField,
			args: args{
				ctx: context.Background(),
				id:  999,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
    test.ResetDatabase(dbUrl, seedPath)
		t.Run(tt.name, func(t *testing.T) {
			r := &SensorRepository{
				querier:                 tt.fields.querier,
				SensorRepositoryMappers: tt.fields.SensorRepositoryMappers,
			}
			got, err := r.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
				return
			}
			assertSensor(t, got, tt.want)
		})
	}
}

func TestSensorRepository_GetStatusByID(t *testing.T) {
	type fields struct {
		querier                 sqlc.Querier
		SensorRepositoryMappers SensorRepositoryMappers
	}
	type args struct {
		ctx context.Context
		id  int32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entities.SensorStatus
		wantErr bool
	}{
		{
			name:   "SensorRepository.GetStatusByID should return online sensor status by id 1",
			fields: defaultField,
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: utils.P(entities.SensorStatusOnline),
		},
		{
			name:   "SensorRepository.GetStatusByID should return offline sensor status by id 2",
			fields: defaultField,
			args: args{
				ctx: context.Background(),
				id:  2,
			},
			want: utils.P(entities.SensorStatusOffline),
		},
		{
			name:   "SensorRepository.GetStatusByID should return unknown sensor status by id 3",
			fields: defaultField,
			args: args{
				ctx: context.Background(),
				id:  3,
			},
			want: utils.P(entities.SensorStatusUnknown),
		},
		{
			name:   "SensorRepository.GetStatusByID should return error when id not found",
			fields: defaultField,
			args: args{
				ctx: context.Background(),
				id:  999,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
    test.ResetDatabase(dbUrl, seedPath)
		t.Run(tt.name, func(t *testing.T) {
			r := &SensorRepository{
				querier:                 tt.fields.querier,
				SensorRepositoryMappers: tt.fields.SensorRepositoryMappers,
			}
			got, err := r.GetStatusByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestSensorRepository_GetSensorByStatus(t *testing.T) {
	type fields struct {
		querier                 sqlc.Querier
		SensorRepositoryMappers SensorRepositoryMappers
	}
	type args struct {
		ctx    context.Context
		status *entities.SensorStatus
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*entities.Sensor
		wantErr bool
	}{
		{
			name:   "SensorRepository.GetSensorByStatus should return all sensors with status online",
			fields: defaultField,
			args: args{
				ctx:    context.Background(),
				status: utils.P(entities.SensorStatusOnline),
			},
			want: []*entities.Sensor{
				{
					ID:     1,
					Status: entities.SensorStatusOnline,
				},
				{
					ID:     4,
					Status: entities.SensorStatusOnline,
				},
			},
		},
		{
			name:   "SensorRepository.GetSensorByStatus should return all sensors with status offline",
			fields: defaultField,
			args: args{
				ctx:    context.Background(),
				status: utils.P(entities.SensorStatusOffline),
			},
			want: []*entities.Sensor{
				{
					ID:     2,
					Status: entities.SensorStatusOffline,
				},
			},
		},
		{
			name:   "SensorRepository.GetSensorByStatus should return all sensors with status unknown",
			fields: defaultField,
			args: args{
				ctx:    context.Background(),
				status: utils.P(entities.SensorStatusUnknown),
			},
			want: []*entities.Sensor{
				{
					ID:     3,
					Status: entities.SensorStatusUnknown,
				},
			},
		},
	}
	for _, tt := range tests {
    test.ResetDatabase(dbUrl, seedPath)
		t.Run(tt.name, func(t *testing.T) {
			r := &SensorRepository{
				querier:                 tt.fields.querier,
				SensorRepositoryMappers: tt.fields.SensorRepositoryMappers,
			}
			got, err := r.GetSensorByStatus(tt.args.ctx, tt.args.status)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Len(t, got, len(tt.want))
			for i := range got {
				assertSensor(t, got[i], tt.want[i])
			}
		})
	}
}

func TestSensorRepository_GetSensorDataByID(t *testing.T) {
	type fields struct {
		querier                 sqlc.Querier
		SensorRepositoryMappers SensorRepositoryMappers
	}
	type args struct {
		ctx context.Context
		id  int32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*entities.SensorData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
    test.ResetDatabase(dbUrl, seedPath)
		t.Run(tt.name, func(t *testing.T) {
			r := &SensorRepository{
				querier:                 tt.fields.querier,
				SensorRepositoryMappers: tt.fields.SensorRepositoryMappers,
			}
			got, err := r.GetSensorDataByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("SensorRepository.GetSensorDataByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SensorRepository.GetSensorDataByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSensorRepository_InsertSensorData(t *testing.T) {
	type fields struct {
		querier                 sqlc.Querier
		SensorRepositoryMappers SensorRepositoryMappers
	}
	type args struct {
		ctx  context.Context
		data []*entities.SensorData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*entities.SensorData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
    test.ResetDatabase(dbUrl, seedPath)
		t.Run(tt.name, func(t *testing.T) {
			r := &SensorRepository{
				querier:                 tt.fields.querier,
				SensorRepositoryMappers: tt.fields.SensorRepositoryMappers,
			}
			got, err := r.InsertSensorData(tt.args.ctx, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("SensorRepository.InsertSensorData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SensorRepository.InsertSensorData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSensorRepository_Create(t *testing.T) {
	type fields struct {
		querier                 sqlc.Querier
		SensorRepositoryMappers SensorRepositoryMappers
	}
	type args struct {
		ctx    context.Context
		sensor *entities.Sensor
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entities.Sensor
		wantErr bool
	}{
		{
			name:   "SensorRepository.Create should create new sensor",
			fields: defaultField,
			args: args{
				ctx: context.Background(),
				sensor: &entities.Sensor{
					Status: entities.SensorStatusOnline,
				},
			},
			want: &entities.Sensor{
				ID:     5,
				Status: entities.SensorStatusOnline,
			},
		},
	}
	for _, tt := range tests {
    test.ResetDatabase(dbUrl, seedPath)
		t.Run(tt.name, func(t *testing.T) {
			r := &SensorRepository{
				querier:                 tt.fields.querier,
				SensorRepositoryMappers: tt.fields.SensorRepositoryMappers,
			}
			got, err := r.Create(tt.args.ctx, tt.args.sensor)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
				return
			}
			assertSensor(t, got, tt.want)
		})
	}
}

func TestSensorRepository_Update(t *testing.T) {
	type fields struct {
		querier                 sqlc.Querier
		SensorRepositoryMappers SensorRepositoryMappers
	}
	type args struct {
		ctx context.Context
		s   *entities.Sensor
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entities.Sensor
		wantErr bool
	}{
		{
			name:   "SensorRepository.Update should update sensor",
			fields: defaultField,
			args: args{
				ctx: context.Background(),
				s: &entities.Sensor{
					ID:     1,
					Status: entities.SensorStatusOffline,
				},
			},
			want: &entities.Sensor{
				ID:     1,
				Status: entities.SensorStatusOffline,
			},
		},
	}
	for _, tt := range tests {
    test.ResetDatabase(dbUrl, seedPath)
		t.Run(tt.name, func(t *testing.T) {
			r := &SensorRepository{
				querier:                 tt.fields.querier,
				SensorRepositoryMappers: tt.fields.SensorRepositoryMappers,
			}
			got, err := r.Update(tt.args.ctx, tt.args.s)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
				return
			}
			assertSensor(t, got, tt.want)
		})
	}
}

func TestSensorRepository_Delete(t *testing.T) {
	type fields struct {
		querier                 sqlc.Querier
		SensorRepositoryMappers SensorRepositoryMappers
	}
	type args struct {
		ctx context.Context
		id  int32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "SensorRepository.Delete should delete sensor",
			fields: defaultField,
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
    test.ResetDatabase(dbUrl, seedPath)
		t.Run(tt.name, func(t *testing.T) {
			r := &SensorRepository{
				querier:                 tt.fields.querier,
				SensorRepositoryMappers: tt.fields.SensorRepositoryMappers,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
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
	assert.Equal(t, got.ID, want.ID)
	assert.Equal(t, got.Status, want.Status)
}
