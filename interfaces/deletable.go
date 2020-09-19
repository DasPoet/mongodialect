package interfaces

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

// Deletable is the interface that wraps functions for deleting documents from
// a collection in a MongoDB database.
type Deletable interface {

	// Delete deletes at most one document from a collection that matches a given filter.
	Delete(ctx context.Context, f Filter) (*mongo.DeleteResult, error)

	// DeleteMany deletes all documents from a collection that match a given filter.
	DeleteMany(ctx context.Context, f Filter) (*mongo.DeleteResult, error)

	// DeleteByID deletes at most one document from a collection that has a given ID.
	DeleteByID(ctx context.Context, id interface{}) (*mongo.DeleteResult, error)
}
