package interfaces

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

// Insertable is the interface that wraps the Insert function which inserts
// objects into a MongoDB database.
type Insertable interface {

	// Insert takes an interface and decodes it into an object, which is then inserted
	// into a collection.
	Insert(ctx context.Context, obj interface{}) (*mongo.InsertOneResult, error)

	// InsertMany takes a variadic number of interfaces and decodes them into objects,
	// which are then inserted into a collection.
	InsertMany(ctx context.Context, obj ...interface{}) (*mongo.InsertManyResult, error)
}
