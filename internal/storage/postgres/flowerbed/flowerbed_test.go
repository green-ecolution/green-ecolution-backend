package flowerbed

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/flowerbed/mapper"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/flowerbed/mapper/generated"
	imgMapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/image/mapper"
	imgMapperImpl "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/image/mapper/generated"
	sensorMapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor/mapper"
	sensorMapperImpl "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/test"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

var (
	seedPath     string
	dbURL        string
	querier      *sqlc.Queries
	mapperRepo   FlowerbedMappers
	defaultField struct {
		querier          *sqlc.Queries
		FlowerbedMappers FlowerbedMappers
	}
)

func TestMain(m *testing.M) {
	rootDir := utils.RootDir()
	seedPath = fmt.Sprintf("%s/internal/storage/postgres/test/seed/flowerbed", rootDir)
	close, url, err := test.SetupPostgresContainer(seedPath)
	if err != nil {
		slog.Error("Error setting up postgres container", "error", err)
		os.Exit(1)
	}
	defer close()

	dbURL = *url
	db, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		os.Exit(1)
	}
	querier = sqlc.New(db)
	mapperRepo = FlowerbedMappers{
		mapper:       &generated.InternalFlowerbedRepoMapperImpl{},
		sensorMapper: &sensorMapperImpl.InternalSensorRepoMapperImpl{},
		imgMapper:    &imgMapperImpl.InternalImageRepoMapperImpl{},
	}
	defaultField = struct {
		querier          *sqlc.Queries
		FlowerbedMappers FlowerbedMappers
	}{
		querier:          querier,
		FlowerbedMappers: mapperRepo,
	}
	os.Exit(m.Run())
}

func TestNewFlowerbedMappers(t *testing.T) {
	type args struct {
		fMapper mapper.InternalFlowerbedRepoMapper
		iMapper imgMapper.InternalImageRepoMapper
		sMapper sensorMapper.InternalSensorRepoMapper
	}
	tests := []struct {
		name string
		args args
		want FlowerbedMappers
	}{
		{
			name: "Test NewFlowerbedMappers",
			args: args{
				fMapper: &generated.InternalFlowerbedRepoMapperImpl{},
				iMapper: &imgMapperImpl.InternalImageRepoMapperImpl{},
				sMapper: &sensorMapperImpl.InternalSensorRepoMapperImpl{},
			},
			want: FlowerbedMappers{
				mapper:       &generated.InternalFlowerbedRepoMapperImpl{},
				imgMapper:    &imgMapperImpl.InternalImageRepoMapperImpl{},
				sensorMapper: &sensorMapperImpl.InternalSensorRepoMapperImpl{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFlowerbedMappers(tt.args.fMapper, tt.args.iMapper, tt.args.sMapper)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewFlowerbedRepository(t *testing.T) {
	type args struct {
		querier *sqlc.Queries
		mappers FlowerbedMappers
	}
	tests := []struct {
		name string
		args args
		want storage.FlowerbedRepository
	}{
		{
			name: "Test NewFlowerbedRepository",
			args: args{
				querier: querier,
				mappers: mapperRepo,
			},
			want: &FlowerbedRepository{
				querier:          querier,
				FlowerbedMappers: mapperRepo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFlowerbedRepository(tt.args.querier, tt.args.mappers)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFlowerbedRepository_GetAll(t *testing.T) {
	type fields struct {
		querier *sqlc.Queries
		FlowerbedMappers FlowerbedMappers
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []*entities.Flowerbed
		wantErr   bool
		runBefore func()
	}{
		{
			name:      "FlowebedRepository.GetAll() should return all flowerbeds",
			fields:    defaultField,
			runBefore: resetDatabase,
			args: args{
				ctx: context.Background(),
			},
			want: []*entities.Flowerbed{
				{
					ID: 1,
					Sensor: &entities.Sensor{
						ID:     1,
						Status: entities.SensorStatusOnline,
					},
					Size:           10,
					Description:    "Blumenbeet am Strand",
					NumberOfPlants: 5,
					MoistureLevel:  0.75,
					Region:         "Mürwik",
					Address:        "Solitüde Strand",
					Archived:       false,
					Latitude:       54.820940,
					Longitude:      9.489022,
					Images: []*entities.Image{
						{
							ID:       1,
							URL:      "https://avatars.githubusercontent.com/u/165842746?s=96&v=4",
							Filename: nil,
							MimeType: nil,
						},
						{
							ID:       2,
							URL:      "https://app.dev.green-ecolution.de/api/v1/images/flowerbed.png",
							Filename: utils.P("flowerbed.png"),
							MimeType: utils.P("image/png"),
						},
					},
				},
				{
					ID: 2,
					Sensor: &entities.Sensor{
						ID:     2,
						Status: entities.SensorStatusOffline,
					},
					Size:           11,
					Description:    "Blumenbeet beim Sankt-Jürgen-Platz",
					NumberOfPlants: 5,
					MoistureLevel:  0.5,
					Region:         "Jürgensby",
					Address:        "Ulmenstraße",
					Archived:       false,
					Latitude:       54.78805731048199,
					Longitude:      9.44400186680097,
					Images:         []*entities.Image{},
				},
			},
		},
	}
	for _, tt := range tests {
		if tt.runBefore != nil {
			tt.runBefore()
		}
		t.Run(tt.name, func(t *testing.T) {
			r := &FlowerbedRepository{
				querier:          tt.fields.querier,
				FlowerbedMappers: tt.fields.FlowerbedMappers,
			}
			got, err := r.GetAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Len(t, got, len(tt.want))
			for i := range got {
				assertFlowerbed(t, got[i], tt.want[i])
			}
		})
	}
}

func TestFlowerbedRepository_GetByID(t *testing.T) {
	type fields struct {
		querier *sqlc.Queries
		FlowerbedMappers FlowerbedMappers
	}
	type args struct {
		ctx context.Context
		id  int32
	}
	tests := []struct {
		name      string
		fields    fields
		runBefore func()
		args      args
		want      *entities.Flowerbed
		wantErr   bool
	}{
		{
			name:      "FlowerbedRepository.GetByID() should return a flowerbed by ID",
			fields:    defaultField,
			runBefore: resetDatabase,
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: &entities.Flowerbed{
				ID: 1,
				Sensor: &entities.Sensor{
					ID:     1,
					Status: entities.SensorStatusOnline,
				},
				Images: []*entities.Image{
					{
						ID:       1,
						URL:      "https://avatars.githubusercontent.com/u/165842746?s=96&v=4",
						Filename: nil,
						MimeType: nil,
					},
					{
						ID:       2,
						URL:      "https://app.dev.green-ecolution.de/api/v1/images/flowerbed.png",
						Filename: utils.P("flowerbed.png"),
						MimeType: utils.P("image/png"),
					},
				},
				Size:           10,
				Description:    "Blumenbeet am Strand",
				NumberOfPlants: 5,
				MoistureLevel:  0.75,
				Region:         "Mürwik",
				Address:        "Solitüde Strand",
				Archived:       false,
				Latitude:       54.820940,
				Longitude:      9.489022,
			},
		},
		{
			name:   "FlowerbedRepository.GetByID() should return an error when flowerbed does not exist",
			fields: defaultField,
			args: args{
				ctx: context.Background(),
				id:  100,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if tt.runBefore != nil {
			tt.runBefore()
		}
		t.Run(tt.name, func(t *testing.T) {
			r := &FlowerbedRepository{
				querier:          tt.fields.querier,
				FlowerbedMappers: tt.fields.FlowerbedMappers,
			}
			got, err := r.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
				return
			}
			assertFlowerbed(t, got, tt.want)
		})
	}
}

func TestFlowerbedRepository_GetAllImagesByID(t *testing.T) {
	type fields struct {
		querier *sqlc.Queries
		FlowerbedMappers FlowerbedMappers
	}
	type args struct {
		ctx         context.Context
		flowerbedID int32
	}
	tests := []struct {
		name      string
		fields    fields
		runBefore func()
		args      args
		want      []*entities.Image
		wantErr   bool
	}{
		{
			name:      "FlowerbedRepository.GetAllImagesByID() should return all images by flowerbed ID",
			fields:    defaultField,
			runBefore: resetDatabase,
			args: args{
				ctx:         context.Background(),
				flowerbedID: 1,
			},
			want: []*entities.Image{
				{
					ID:       1,
					URL:      "https://avatars.githubusercontent.com/u/165842746?s=96&v=4",
					Filename: nil,
					MimeType: nil,
				},
				{
					ID:       2,
					URL:      "https://app.dev.green-ecolution.de/api/v1/images/flowerbed.png",
					Filename: utils.P("flowerbed.png"),
					MimeType: utils.P("image/png"),
				},
			},
		},
		{
			name:   "FlowerbedRepository.GetAllImagesByID() should return an empty list when flowerbed does not have images",
			fields: defaultField,
			args: args{
				ctx:         context.Background(),
				flowerbedID: 2,
			},
			want: []*entities.Image{},
		},
	}
	for _, tt := range tests {
		if tt.runBefore != nil {
			tt.runBefore()
		}
		t.Run(tt.name, func(t *testing.T) {
			r := &FlowerbedRepository{
				querier:          tt.fields.querier,
				FlowerbedMappers: tt.fields.FlowerbedMappers,
			}
			got, err := r.GetAllImagesByID(tt.args.ctx, tt.args.flowerbedID)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Len(t, got, len(tt.want))
			for i := range got {
				assertImage(t, got[i], tt.want[i])
			}
		})
	}
}

func TestFlowerbedRepository_GetSensorByFlowerbedID(t *testing.T) {
	type fields struct {
		querier *sqlc.Queries
		FlowerbedMappers FlowerbedMappers
	}
	type args struct {
		ctx         context.Context
		flowerbedID int32
	}
	tests := []struct {
		name      string
		runBefore func()
		fields    fields
		args      args
		want      *entities.Sensor
		wantErr   bool
	}{
		{
			name:      "FlowerbedRepository.GetSensorByFlowerbedID() should return a sensor by flowerbed ID",
			runBefore: resetDatabase,
			fields:    defaultField,
			args: args{
				ctx:         context.Background(),
				flowerbedID: 1,
			},
			want: &entities.Sensor{
				ID:     1,
				Status: entities.SensorStatusOnline,
			},
		},
		{
			name:   "FlowerbedRepository.GetSensorByFlowerbedID() should return an error when sensor does not exist",
			fields: defaultField,
			args: args{
				ctx:         context.Background(),
				flowerbedID: 100,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if tt.runBefore != nil {
			tt.runBefore()
		}
		t.Run(tt.name, func(t *testing.T) {
			r := &FlowerbedRepository{
				querier:          tt.fields.querier,
				FlowerbedMappers: tt.fields.FlowerbedMappers,
			}
			got, err := r.GetSensorByFlowerbedID(tt.args.ctx, tt.args.flowerbedID)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
				return
			}
			assertSensor(t, got, tt.want)
		})
	}
}

func TestFlowerbedRepository_Create(t *testing.T) {
	type fields struct {
		querier *sqlc.Queries
		FlowerbedMappers FlowerbedMappers
	}
	type args struct {
		ctx context.Context
		f   *entities.CreateFlowerbed
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		runBefore func()
		runAfter  func()
		want      *entities.Flowerbed
		wantErr   bool
		errType   error
	}{
		{
			name:   "FlowerbedRepository.Create() should create a flowerbed",
			fields: defaultField,
			runBefore: func() {
				resetDatabase()
				arg := &sqlc.CreateImageParams{
					Filename: utils.P("flowerbed.png"),
					MimeType: utils.P("image/png"),
					Url:      "https://app.dev.green-ecolution.de/api/v1/images/flowerbed.png",
				}
				_, err := querier.CreateImage(context.Background(), arg)
				if err != nil {
					return
				}
			},
			args: args{
				ctx: context.Background(),
				f: &entities.CreateFlowerbed{
					Size:           10,
					Description:    "Blumenbeet am Strand",
					NumberOfPlants: 5,
					MoistureLevel:  0.75,
					Region:         "Mürwik",
					Address:        "Solitüde Strand",
					Archived:       false,
					Latitude:       54.820940,
					Longitude:      9.489022,
					SensorID:       utils.P(int32(1)),
					ImageIDs:       []int32{3},
				},
			},
			want: &entities.Flowerbed{
				ID:             3,
				Size:           10,
				Description:    "Blumenbeet am Strand",
				NumberOfPlants: 5,
				MoistureLevel:  0.75,
				Region:         "Mürwik",
				Address:        "Solitüde Strand",
				Archived:       false,
				Latitude:       54.820940,
				Longitude:      9.489022,
				Sensor: &entities.Sensor{
					ID:     1,
					Status: entities.SensorStatusOnline,
				},
				Images: []*entities.Image{
					{
						ID:       3,
						URL:      "https://app.dev.green-ecolution.de/api/v1/images/flowerbed.png",
						Filename: utils.P("flowerbed.png"),
						MimeType: utils.P("image/png"),
					},
				},
			},
		},

		{
			name:      "should return an error when sensor does not exist and not create entity",
			fields:    defaultField,
			runBefore: resetDatabase,
			args: args{
				ctx: context.Background(),
				f: &entities.CreateFlowerbed{
					Size:           10,
					Description:    "Blumenbeet am Strand",
					NumberOfPlants: 5,
					MoistureLevel:  0.75,
					Region:         "Mürwik",
					Address:        "Solitüde Strand",
					Archived:       false,
					Latitude:       54.820940,
					Longitude:      9.489022,
					SensorID:       utils.P(int32(100)),
					ImageIDs:       []int32{1, 2},
				},
			},
			want:    nil,
			wantErr: true,
			errType: storage.ErrSensorNotFound,
			runAfter: func() {
				r := &FlowerbedRepository{
					querier:          defaultField.querier,
					FlowerbedMappers: defaultField.FlowerbedMappers,
				}
				row, err := r.GetByID(context.Background(), 3)
				assert.Nil(t, row)
				assert.Error(t, err)
			},
		},

		{
			name:      "should return an error when image does not exist and not create entity",
			fields:    defaultField,
			runBefore: resetDatabase,
			args: args{
				ctx: context.Background(),
				f: &entities.CreateFlowerbed{
					Size:           10,
					Description:    "Blumenbeet am Strand",
					NumberOfPlants: 5,
					MoistureLevel:  0.75,
					Region:         "Mürwik",
					Address:        "Solitüde Strand",
					Archived:       false,
					Latitude:       54.820940,
					Longitude:      9.489022,
					SensorID:       utils.P(int32(1)),
					ImageIDs:       []int32{1, 100},
				},
			},
			want:    nil,
			wantErr: true,
			errType: storage.ErrImageNotFound,
			runAfter: func() {
				r := &FlowerbedRepository{
					querier:          defaultField.querier,
					FlowerbedMappers: defaultField.FlowerbedMappers,
				}
				row, err := r.GetByID(context.Background(), 3)
				assert.Nil(t, row)
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		if tt.runBefore != nil {
			tt.runBefore()
		}
		t.Run(tt.name, func(t *testing.T) {
			r := &FlowerbedRepository{
				querier:          tt.fields.querier,
				FlowerbedMappers: tt.fields.FlowerbedMappers,
			}
			got, err := r.Create(tt.args.ctx, tt.args.f)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.errType)
				return
			}
			assertFlowerbed(t, got, tt.want)

			if tt.runAfter != nil {
				tt.runAfter()
			}
		})
	}
}

func TestFlowerbedRepository_Update(t *testing.T) {
	type fields struct {
		querier *sqlc.Queries
		FlowerbedMappers FlowerbedMappers
	}
	type args struct {
		ctx context.Context
		f   *entities.UpdateFlowerbed
	}
	tests := []struct {
		name      string
		fields    fields
		runBefore func()
    runAfter  func()
		args      args
		want      *entities.Flowerbed
		wantErr   bool
		errType   error
	}{
		{
			name:   "should update all fields of a flowerbed when all fields are provided",
			fields: defaultField,
			runBefore: func() {
				// Create Image with ID 3 and 4
				for i := 3; i <= 4; i++ {
					arg := &sqlc.CreateImageParams{
						Filename: utils.P("flowerbed.png"),
						MimeType: utils.P("image/png"),
						Url:      "https://app.dev.green-ecolution.de/api/v1/images/flowerbed.png",
					}
					_, err := querier.CreateImage(context.Background(), arg)
					if err != nil {
						return
					}
				}
			},
      runAfter: func() {
        for i := 3; i <= 4; i++ {
          if err := querier.DeleteImage(context.Background(), int32(i)); err != nil {
            return
          }
        }
      },
			args: args{
				ctx: context.Background(),
				f: &entities.UpdateFlowerbed{
					ID:             1,
					Size:           utils.P(10.3),
					Description:    utils.P("Neues Blumenbeet am Strand"),
					NumberOfPlants: utils.P(int32(129)),
					MoistureLevel:  utils.P(0.38),
					Region:         utils.P("Neustadt"),
					Address:        utils.P("Wo auch immer"),
					Archived:       utils.P(false),
					Latitude:       utils.P(54.820938),
					Longitude:      utils.P(9.489392),
					SensorID:       utils.P(int32(2)),
					ImageIDs:       []int32{3, 4},
				},
			},
			want: &entities.Flowerbed{
				ID:             1,
        Size:           10.3,
        Description:    "Neues Blumenbeet am Strand",
        NumberOfPlants: 129,
        MoistureLevel:  0.38,
        Region:         "Neustadt",
        Address:        "Wo auch immer",
        Archived:       false,
        Latitude:       54.820938,
        Longitude:      9.489392,
				Sensor: &entities.Sensor{
					ID:     2,
					Status: entities.SensorStatusOffline,
				},
				Images: []*entities.Image{
					{
						ID:       3,
						URL:      "https://app.dev.green-ecolution.de/api/v1/images/flowerbed.png",
						Filename: utils.P("flowerbed.png"),
						MimeType: utils.P("image/png"),
					},
					{
						ID:       4,
						URL:      "https://app.dev.green-ecolution.de/api/v1/images/flowerbed.png",
						Filename: utils.P("flowerbed.png"),
						MimeType: utils.P("image/png"),
					},
				},
			},
		},

		{
			name:      "should update only the fields that are provided",
			fields:    defaultField,
			args: args{
				ctx: context.Background(),
				f: &entities.UpdateFlowerbed{
					ID:             1,
					Size:           utils.P(12.0),
					Description:    utils.P("FooBar"),
					NumberOfPlants: nil,
					MoistureLevel:  nil,
					Region:         nil,
					Address:        nil,
					Archived:       nil,
					Latitude:       nil,
					Longitude:      nil,
					SensorID:       nil,
					ImageIDs:       nil,
				},
			},
			want: &entities.Flowerbed{
				ID:             1,
				Size:           12.0,
				Description:    "FooBar",
				NumberOfPlants: 5,
				MoistureLevel:  0.75,
				Region:         "Mürwik",
				Address:        "Solitüde Strand",
				Archived:       false,
				Latitude:       54.820940,
				Longitude:      9.489022,
				Sensor: &entities.Sensor{
					ID:     1,
					Status: entities.SensorStatusOnline,
				},
				Images: []*entities.Image{
					{
						ID:       1,
						URL:      "https://avatars.githubusercontent.com/u/165842746?s=96&v=4",
						Filename: nil,
						MimeType: nil,
					},
					{
						ID:       2,
						URL:      "https://app.dev.green-ecolution.de/api/v1/images/flowerbed.png",
						Filename: utils.P("flowerbed.png"),
						MimeType: utils.P("image/png"),
					},
				},
			},
		},

		{
			name:      "should update images when images IDs is provided",
			fields:    defaultField,
			args: args{
				ctx: context.Background(),
				f: &entities.UpdateFlowerbed{
					ID:             1,
					Size:           utils.P(12.0),
					Description:    utils.P("Blumenbeet in der Stadt"),
					NumberOfPlants: utils.P(int32(5)),
					MoistureLevel:  utils.P(0.75),
					Region:         utils.P("Mürwik"),
					Address:        utils.P("Solitüde Strand"),
					Archived:       utils.P(false),
					Latitude:       utils.P(54.820940),
					Longitude:      utils.P(9.489022),
					SensorID:       utils.P(int32(1)),
					ImageIDs:       []int32{1}, // should remove image with ID 2
				},
			},
			want: &entities.Flowerbed{
				ID:             1,
				Size:           12.0,
				Description:    "Blumenbeet in der Stadt",
				NumberOfPlants: 5,
				MoistureLevel:  0.75,
				Region:         "Mürwik",
				Address:        "Solitüde Strand",
				Archived:       false,
				Latitude:       54.820940,
				Longitude:      9.489022,
				Sensor: &entities.Sensor{
					ID:     1,
					Status: entities.SensorStatusOnline,
				},
				Images: []*entities.Image{
					{
						ID:       1,
						URL:      "https://avatars.githubusercontent.com/u/165842746?s=96&v=4",
						Filename: nil,
						MimeType: nil,
					},
				},
			},
		},

		{
			name:      "should update sensor when sensor ID is provided",
			fields:    defaultField,
			args: args{
				ctx: context.Background(),
				f: &entities.UpdateFlowerbed{
					ID:       1,
					SensorID: utils.P(int32(2)),
				},
			},
			want: &entities.Flowerbed{
				ID:             1,
				Size:           12.0,
				Description:    "Blumenbeet in der Stadt",
				NumberOfPlants: 5,
				MoistureLevel:  0.75,
				Region:         "Mürwik",
				Address:        "Solitüde Strand",
				Archived:       false,
				Latitude:       54.820940,
				Longitude:      9.489022,
				Sensor: &entities.Sensor{
					ID:     2,
					Status: entities.SensorStatusOffline,
				},
				Images: []*entities.Image{
					{
						ID:       1,
						URL:      "https://avatars.githubusercontent.com/u/165842746?s=96&v=4",
						Filename: nil,
						MimeType: nil,
					},
				},
			},
		},

		{
			name:      "should return an error when flowerbed does not exist",
			fields:    defaultField,
			args: args{
				ctx: context.Background(),
				f: &entities.UpdateFlowerbed{
					ID: 100,
				},
			},
			want:    nil,
			wantErr: true,
			errType: storage.ErrIDNotFound,
		},

		{
			name:      "should return an error when sensor does not exist",
			fields:    defaultField,
			args: args{
				ctx: context.Background(),
				f: &entities.UpdateFlowerbed{
					ID:       1,
					SensorID: utils.P(int32(100)),
				},
			},
			want:    nil,
			wantErr: true,
			errType: storage.ErrSensorNotFound,
		},

		{
			name:      "should return an error when image does not exist",
			fields:    defaultField,
			args: args{
				ctx: context.Background(),
				f: &entities.UpdateFlowerbed{
					ID:       1,
					ImageIDs: []int32{100},
				},
			},
			want:    nil,
			wantErr: true,
			errType: storage.ErrImageNotFound,
		},
	}
	for _, tt := range tests {
		if tt.runBefore != nil {
			tt.runBefore()
		}

		t.Run(tt.name, func(t *testing.T) {
			cleanup, err := setupTx(tt.args.ctx)
			if err != nil {
				return
			}
			defer cleanup()
			r := &FlowerbedRepository{
				querier:          tt.fields.querier,
				FlowerbedMappers: tt.fields.FlowerbedMappers,
			}
			got, err := r.Update(tt.args.ctx, tt.args.f)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.errType)
				return
			}
			assertFlowerbed(t, got, tt.want)
      if tt.runAfter != nil {
        tt.runAfter()
      }
		})
	}
}

func TestFlowerbedRepository_Delete(t *testing.T) {
	type fields struct {
		querier *sqlc.Queries
		FlowerbedMappers FlowerbedMappers
	}
	type args struct {
		ctx context.Context
		id  int32
	}
	tests := []struct {
		name      string
		fields    fields
		runBefore func()
		args      args
		wantErr   bool
	}{
		{
			name:      "FlowerbedRepository.Delete() should delete a flowerbed",
			fields:    defaultField,
			runBefore: resetDatabase,
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			wantErr: false,
		},
		{
			name:   "FlowerbedRepository.Delete() should not return an error when flowerbed does not exist",
			fields: defaultField,
			args: args{
				ctx: context.Background(),
				id:  100,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		if tt.runBefore != nil {
			tt.runBefore()
		}
		t.Run(tt.name, func(t *testing.T) {
			r := &FlowerbedRepository{
				querier:          tt.fields.querier,
				FlowerbedMappers: tt.fields.FlowerbedMappers,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
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

func resetDatabase() {
	if err := test.ResetDatabase(dbURL, seedPath); err != nil {
		slog.Error("Error resetting database: %v", "error", err)
		os.Exit(1)
	}
}

func setupTx(ctx context.Context) (cleanup func(), err error) {
	db, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	qtx := querier.WithTx(tx)
	querier = qtx

	cleanup = func() {
    log.Println("Rolling back transaction")
		if err := tx.Rollback(ctx); err != nil {
			slog.Error("Error rolling back transaction", "error", err)
		}
	}

	return cleanup, nil
}
