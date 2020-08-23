package interfaces

// Repository is the interface that wraps interfaces providing basic interactions
// with collections in a MongoDB database.
type Repository interface {

	// See Findable
	Findable

	// See Insertable
	Insertable

	// See Updatable
	Updatable

	// See Deletable
	Deletable

	// Type returns the underlying type of the documents in the collection the
	// Repository is able to access.
	Type() interface{}
}