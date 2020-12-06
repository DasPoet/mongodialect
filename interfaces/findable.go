package interfaces

import (
    "context"
)

// Findable is the interface that wraps functions for finding and checking
// the existence of document sin a collection in a MongoDB database.
type Findable interface {

    // Find decodes all documents from a collection that match a given filter, and returns them.
    Find(ctx context.Context, f Filter) ([]interface{}, error)

    // FindByID decodes one document from a collection that has a given ID, and returns it.
    FindByID(ctx context.Context, id interface{}) (interface{}, error)

    // Exists returns whether at least one document matching a given filter exists in a collection.
    Exists(ctx context.Context, f Filter) (bool, error)

    // ExistsByID returns whether at least one document having a given ID exists in a collection.
    ExistsByID(ctx context.Context, id interface{}) (bool, error)
}
