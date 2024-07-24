package tree

import (
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities/tree"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// goverter:converter
// goverter:extend PrimitiveIDToString StringToPrimitiveID
type TreeMongoMapper interface {
	FromEntity(src *TreeEntity) *domain.Tree
	FromEntityList(src []*TreeEntity) []*domain.Tree

	ToEntity(src *domain.Tree) *TreeEntity
	ToEntityList(src []*domain.Tree) []*TreeEntity
}

func PrimitiveIDToString(id primitive.ObjectID) string {
	return id.Hex()
}

func StringToPrimitiveID(id string) primitive.ObjectID {
	objID, _ := primitive.ObjectIDFromHex(id)
	return objID
}
