package wateringplan

import (
	"context"
	"log/slog"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
)

type WateringPlanRepository struct {
	store *store.Store
	WateringPlanMappers
}

type WateringPlanMappers struct {
	mapper        mapper.InternalWateringPlanRepoMapper
	vehicleMapper mapper.InternalVehicleRepoMapper
	clusterMapper mapper.InternalTreeClusterRepoMapper
}

func NewWateringPlanRepositoryMappers(
	wMapper mapper.InternalWateringPlanRepoMapper,
	vMapper mapper.InternalVehicleRepoMapper,
	tcMapper mapper.InternalTreeClusterRepoMapper,
) WateringPlanMappers {
	return WateringPlanMappers{
		mapper:        wMapper,
		vehicleMapper: vMapper,
		clusterMapper: tcMapper,
	}
}

func NewWateringPlanRepository(s *store.Store, mappers WateringPlanMappers) storage.WateringPlanRepository {
	s.SetEntityType(store.WateringPlan)
	return &WateringPlanRepository{
		store:               s,
		WateringPlanMappers: mappers,
	}
}

func WithDate(date time.Time) entities.EntityFunc[entities.WateringPlan] {
	return func(wp *entities.WateringPlan) {
		slog.Debug("updating date", "date", date)
		wp.Date = date
	}
}

func WithDescription(description string) entities.EntityFunc[entities.WateringPlan] {
	return func(wp *entities.WateringPlan) {
		slog.Debug("updating description", "description", description)
		wp.Description = description
	}
}

func WithWateringPlanStatus(status entities.WateringPlanStatus) entities.EntityFunc[entities.WateringPlan] {
	return func(wp *entities.WateringPlan) {
		slog.Debug("updating watering plan status", "watering plan status", status)
		wp.Status = status
	}
}

func WithDistance(distance *float64) entities.EntityFunc[entities.WateringPlan] {
	return func(wp *entities.WateringPlan) {
		slog.Debug("updating distance", "watering distance", distance)
		wp.Distance = distance
	}
}

func WithTotalWaterRequired(totalWaterRequired *float64) entities.EntityFunc[entities.WateringPlan] {
	return func(wp *entities.WateringPlan) {
		slog.Debug("updating total water required", "total water required", totalWaterRequired)
		wp.TotalWaterRequired = totalWaterRequired
	}
}

func WithUsers(users []*entities.User) entities.EntityFunc[entities.WateringPlan] {
	return func(wp *entities.WateringPlan) {
		slog.Debug("updating users", "users", users)
		wp.Users = users
	}
}

func WithTreecluster(treecluster []*entities.TreeCluster) entities.EntityFunc[entities.WateringPlan] {
	return func(wp *entities.WateringPlan) {
		slog.Debug("updating tree cluster", "tree cluster", treecluster)
		wp.Treecluster = treecluster
	}
}

func WithTransporter(transporter *entities.Vehicle) entities.EntityFunc[entities.WateringPlan] {
	return func(wp *entities.WateringPlan) {
		slog.Debug("updating transporter", "transporter", transporter)
		wp.Transporter = transporter
	}
}

func WithTrailer(trailer *entities.Vehicle) entities.EntityFunc[entities.WateringPlan] {
	return func(wp *entities.WateringPlan) {
		slog.Debug("updating trailer", "trailer", trailer)
		wp.Trailer = trailer
	}
}

func (w *WateringPlanRepository) Delete(ctx context.Context, id int32) error {
	_, err := w.store.DeleteWateringPlan(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
