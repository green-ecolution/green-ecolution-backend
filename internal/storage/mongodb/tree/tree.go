package tree

import (
	"context"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities/tree"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/mongodb/entities/tree"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/mongodb/entities/tree/generated"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TreeRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
	mapper     tree.TreeMongoMapper
}

func NewTreeRepository(client *mongo.Client, collection *mongo.Collection) *TreeRepository {
	return &TreeRepository{client: client, collection: collection, mapper: &generated.TreeMongoMapperImpl{}}
}

func (r *TreeRepository) Insert(ctx context.Context, data *domain.Tree) error {
	_, err := r.collection.InsertOne(ctx, data)
	if err != nil {
		return storage.ErrMongoCannotUpsertData
	}

	return nil
}

func (r *TreeRepository) Get(ctx context.Context, id string) (*domain.Tree, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var data tree.TreeEntity
	err = r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&data)
	if err != nil {
		return nil, storage.ErrMongoDataNotFound
	}

	domainData := r.mapper.FromEntity(&data)
	return domainData, nil
}

func (r *TreeRepository) GetAll(ctx context.Context) ([]*domain.Tree, error) {
	var data []*tree.TreeEntity
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, storage.ErrMongoDataNotFound
	}
	if err = cursor.All(ctx, &data); err != nil {
		return nil, storage.ErrMongoDataNotFound
	}

	domainData := r.mapper.FromEntityList(data)
	return domainData, nil
}
