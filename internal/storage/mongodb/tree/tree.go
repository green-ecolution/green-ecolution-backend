package tree

import (
	"context"

	"github.com/SmartCityFlensburg/green-space-management/internal/entities/tree"
	"github.com/SmartCityFlensburg/green-space-management/internal/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TreeRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewTreeRepository(client *mongo.Client, collection *mongo.Collection) *TreeRepository {
	return &TreeRepository{client: client, collection: collection}
}

func (r *TreeRepository) Insert(ctx context.Context, data tree.Tree) error {
	_, err := r.collection.InsertOne(ctx, data)
	if err != nil {
		return storage.ErrMongoCannotUpsertData
	}

	return nil
}

func (r *TreeRepository) Get(ctx context.Context, id string) (*tree.Tree, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var data tree.Tree
	err = r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&data)
	if err != nil {
		return nil, storage.ErrMongoDataNotFound
	}

	return &data, nil
}

func (r *TreeRepository) GetAll(ctx context.Context) ([]tree.Tree, error) {
	var data []tree.Tree
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, storage.ErrMongoDataNotFound
	}
	if err = cursor.All(ctx, &data); err != nil {
		return nil, storage.ErrMongoDataNotFound
	}

	return data, nil
}