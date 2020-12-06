package interfaces

import (
    "context"
    "go.mongodb.org/mongo-driver/mongo"
)

// Updatable is the interface that wraps functions for
// updating documents in a collection in a MongoDB database.
type Updatable interface {

    // Update decodes changes into a value which is subsequently
    // used to update the first document matching f.
    Update(ctx context.Context, f Filter, changes map[string]interface{}) (*mongo.UpdateResult, error)

    // UpdateByID decodes changes into a value which is subsequently
    // used to update the first document having the given id.
    UpdateByID(ctx context.Context, id interface{}, changes map[string]interface{}) (*mongo.UpdateResult, error)
}
