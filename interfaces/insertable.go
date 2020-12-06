package interfaces

import (
    "context"
    "go.mongodb.org/mongo-driver/mongo"
)

// Insertable is the interface that wraps functions
// for inserting documents into a MongoDB database.
type Insertable interface {

    // Insert decodes v into a document and inserts it into a collection.
    Insert(ctx context.Context, v interface{}) (*mongo.InsertOneResult, error)

    // InsertMany decodes every element in v and inserts it into a collection.
    InsertMany(ctx context.Context, v ...interface{}) (*mongo.InsertManyResult, error)
}
