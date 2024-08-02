package sensor

import (
	"context"
	"log/slog"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/mongodb/entities/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/mongodb/entities/sensor/generated"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SensorRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
	mapper     sensor.MqttMongoMapper
}

func NewSensorRepository(client *mongo.Client, collection *mongo.Collection) *SensorRepository {
	return &SensorRepository{client: client, collection: collection, mapper: &generated.MqttMongoMapperImpl{}}
}

func (r *SensorRepository) Insert(ctx context.Context, data *domain.MqttPayload) (*domain.MqttPayload, error) {
	payloadEntity := r.mapper.ToEntity(data)
	entity := &sensor.MqttEntity{
		Data:   *payloadEntity,
		TreeID: "6686f54fd32cf640e8ae6eb1",
	}

	if entity.ID == primitive.NilObjectID {
		objID := primitive.NewObjectID()
		entity.ID = objID
	}
	_, err := r.collection.InsertOne(ctx, data)
	if err != nil {
		return nil, storage.ErrMongoCannotUpsertData
	}

	return data, nil
}

func (r *SensorRepository) Get(ctx context.Context, id string) (*domain.MqttPayload, error) {
	filter := bson.M{"end_device_ids.device_id": id}
	var data sensor.MqttEntity
	err := r.collection.FindOne(ctx, filter).Decode(&data)
	if err != nil {
		return nil, storage.ErrMongoDataNotFound
	}

	domainData := r.mapper.FromEntity(&data)
	return domainData, nil
}

func (r *SensorRepository) GetFirst(ctx context.Context) (*domain.MqttPayload, error) {
	var data sensor.MqttEntity
	if err := r.collection.FindOne(ctx, bson.D{}).Decode(&data); err != nil {
		return nil, storage.ErrMongoDataNotFound
	}

	domainData := r.mapper.FromEntity(&data)
	return domainData, nil
}

func (r *SensorRepository) GetAllByTreeID(ctx context.Context, treeID string) ([]*domain.MqttPayload, error) {
	filter := bson.M{"tree_id": treeID}
	var data []*sensor.MqttEntity
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		slog.Error("Error while getting sensor data", "error", err)
		return nil, storage.ErrMongoDataNotFound
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var d sensor.MqttEntity
		if err := cursor.Decode(&d); err != nil {
			slog.Error("Error while decoding sensor data", "error", err)
			return nil, storage.ErrMongoDataNotFound
		}
		data = append(data, &d)
	}

	domainData := r.mapper.FromEntityList(data)
	return domainData, nil
}

func (r *SensorRepository) GetLastByTreeID(ctx context.Context, treeID string) (*domain.MqttPayload, error) {
	filter := bson.M{"tree_id": treeID}
	opts := options.FindOne().SetSort(bson.D{{Key: "data.received_at", Value: -1}})
	var data sensor.MqttEntity
	err := r.collection.FindOne(ctx, filter, opts).Decode(&data)
	if err != nil {
		return nil, storage.ErrMongoDataNotFound
	}

	domainData := r.mapper.FromEntity(&data)
	return domainData, nil
}
