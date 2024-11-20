package sensor

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

func (r *SensorRepository) Update(ctx context.Context, id string, sFn ...entities.EntityFunc[entities.Sensor]) (*entities.Sensor, error) {
	entity, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, fn := range sFn {
		fn(entity)
	}

	if err := r.updateEntity(ctx, entity); err != nil {
		return nil, r.store.HandleError(err)
	}

	if len(entity.Data) > 0 {
		_, err := r.InsertSensorData(ctx, entity.Data)
		if err != nil {
			return nil, err
		}
	}

	return r.GetByID(ctx, entity.ID)
}

func (r *SensorRepository) updateEntity(ctx context.Context, sensor *entities.Sensor) error {
	params := sqlc.UpdateSensorParams{
		ID:     sensor.ID,
		Status: sqlc.SensorStatus(sensor.Status),
	}

	return r.store.UpdateSensor(ctx, &params)
}
