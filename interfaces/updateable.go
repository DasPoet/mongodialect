package interfaces

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

// Updatable is the interface that wraps functions for updating documents in a
// collection in a MongoDB database.
type Updatable interface {

	// Update takes a filter, and a map of changes. It decodes the changes into
	// an object which is then used to update the first document matching the
	// provided filter.
	Update(ctx context.Context, f Filter, changes map[string]interface{}) (*mongo.UpdateResult, error)

	// Update by ID takes an ID, and a map of changes. It decodes the changes into
	// an object which is then used to update the first document having the provided ID.
	UpdateByID(ctx context.Context, id uuid.UUID, changes map[string]interface{}) (*mongo.UpdateResult, error)
}
