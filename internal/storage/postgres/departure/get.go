package depature

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
)

func (d *DepartureRepository) GetAll(ctx context.Context) ([]*entities.Departure, error) {
	rows, err := d.store.GetAllDepartures(ctx)
	if err != nil {
		return nil, d.store.HandleError(err)
	}

	return d.mapper.FromSqlList(rows), nil
}

func (d *DepartureRepository) GetByID(ctx context.Context, id int32) (*entities.Departure, error) {
	row, err := d.store.GetDepartureByID(ctx, id)
	if err != nil {
		return nil, d.store.HandleError(err)
	}

	return d.mapper.FromSql(row), nil
}
