package sensor

import (
	"context"
	"encoding/json"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor/mapper"
	"github.com/pkg/errors"
)

type SensorRepository struct {
	querier sqlc.Querier
	SensorRepositoryMappers
}

type SensorRepositoryMappers struct {
	mapper mapper.InternalSensorRepoMapper
}

func NewSensorRepositoryMappers(sMapper mapper.InternalSensorRepoMapper) SensorRepositoryMappers {
	return SensorRepositoryMappers{
		mapper: sMapper,
	}
}

func NewSensorRepository(querier sqlc.Querier, mappers SensorRepositoryMappers) storage.SensorRepository {
	return &SensorRepository{
		querier:                 querier,
		SensorRepositoryMappers: mappers,
	}
}

func (r *SensorRepository) GetAll(ctx context.Context) ([]*entities.Sensor, error) {
	rows, err := r.querier.GetAllSensors(ctx)
	if err != nil {
		return nil, err
	}

	return r.mapper.FromSqlList(rows), nil
}

func (r *SensorRepository) GetByID(ctx context.Context, id int32) (*entities.Sensor, error) {
	row, err := r.querier.GetSensorByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.mapper.FromSql(row), nil
}

func (r *SensorRepository) GetStatusByID(ctx context.Context, id int32) (*entities.SensorStatus, error) {
	sensor, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &sensor.Status, nil
}

func (r *SensorRepository) GetSensorByStatus(ctx context.Context, status *entities.SensorStatus) ([]*entities.Sensor, error) {
	row, err := r.querier.GetSensorByStatus(ctx, sqlc.SensorStatus(*status))
	if err != nil {
		return nil, err
	}

	return r.mapper.FromSqlList(row), nil
}

func (r *SensorRepository) GetSensorDataByID(ctx context.Context, id int32) ([]*entities.SensorData, error) {
	rows, err := r.querier.GetSensorDataBySensorID(ctx, id)
	if err != nil {
		return nil, err
	}

	domainData := make([]*entities.SensorData, len(rows))

	for i, row := range rows {
		domainData[i] = r.mapper.FromSqlSensorData(row)
		data, err := mapper.MapSensorData(row.Data)
		if err != nil {
			return nil, errors.Wrap(err, "failed to map sensor data")
		}
		domainData[i].Data = data
	}

	return domainData, nil
}

func (r *SensorRepository) InsertSensorData(ctx context.Context, data []*entities.SensorData) ([]*entities.SensorData, error) {
	for _, d := range data {
		mqttData := r.mapper.FromDomainSensorData(d.Data)
		raw, err := json.Marshal(mqttData)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal mqtt data")
		}

		params := &sqlc.InsertSensorDataParams{
			SensorID: d.ID,
			Data:     raw,
		}

		err = r.querier.InsertSensorData(ctx, params)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (r *SensorRepository) Create(ctx context.Context, sensor *entities.Sensor) (*entities.Sensor, error) {
	id, err := r.querier.CreateSensor(ctx, sqlc.SensorStatus(sensor.Status))
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *SensorRepository) Update(ctx context.Context, s *entities.Sensor) (*entities.Sensor, error) {
	params := &sqlc.UpdateSensorParams{
		ID:     s.ID,
		Status: sqlc.SensorStatus(s.Status),
	}
	err := r.querier.UpdateSensor(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, s.ID)
}

func (r *SensorRepository) Delete(ctx context.Context, id int32) error {
	return r.querier.DeleteSensor(ctx, id)
}
