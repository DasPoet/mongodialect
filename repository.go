package mongodialect

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/informatik-q2/mongodialect/interfaces"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"strings"
)

var (
	// ErrorDriverNil occurs when a given driver is nil
	ErrorDriverNil = errors.New("driver must no be nil")

	// ErrorCollectionEmpty occurs when the name of a given collection is empty
	ErrorCollectionEmpty = errors.New("collection must not be ''")

	// ErrorObjectNotFound occurs when a lookup does not yield a result
	ErrorObjectNotFound = errors.New("object was not found")

	// ErrorMultipleMatches occurs when a lookup using a given ID yields more than one result
	ErrorMultipleMatches = errors.New("multiple matches for ID")
)

// A Repository wraps the Driver and provides functionality for performing
// operations on the collections of a MongoDB database. It can only access
// one collection at a time, though.
//
// The Repository implements the Repository interface. It therefore offers the following functionality:
// See: Find, Insert, Update, Delete
// To allow for generic access to collections containing arbitrary data, the
// Repository must be provided with the underlying type of the data structure
// contained in the specific collection to access.
//
type Repository struct {
	baseType   interface{} // the type of the data structure stored in the collection; must be a pointer
	idField    string      // the name of the field containing the underlying object's ID
	collection string      // the name of the collection to access
	Driver     *Driver     // a pointer to a Driver which is used to connect to the database
}

// NewRepository returns a new repository upon validating the given base type and driver.
//
// If the provided ID field is an empty string, the default Mongo ID ("_id") is used instead.
//
// It fails if the driver is nil, or the provided base type is not a pointer.
//
// It also fails if the provided collection is an empty string.
//
func NewRepository(baseType interface{}, driver *Driver, collection string, idField string) (*Repository, error) {
	if driver == nil {
		return nil, ErrorDriverNil
	}

	if strings.EqualFold(strings.Trim(collection, " "), "") {
		return nil, ErrorCollectionEmpty
	}

	rt := reflect.TypeOf(baseType)
	kind := rt.Kind()

	if kind != reflect.Ptr {
		return nil, fmt.Errorf("baseType must be a pointer, not '%s'", kind)
	}

	// fallback to Mongo's default ID
	if strings.EqualFold(strings.Trim(idField, " "), "") {
		idField = "_id"
	}

	return &Repository{
		baseType:   baseType,
		idField:    idField,
		collection: collection,
		Driver:     driver,
	}, nil
}

// Type returns a pointer to the type Repository's base type, which has the same
// type as the data stored in the Repository's collection.
func (repository *Repository) Type() interface{} {
	return repository.baseType
}

// Find decodes all documents in the Repository's collection that match a given filter.
//
// It accesses the Repository's collection using the provided filter, and decodes
// the result into the Repository's base type.
//
// It fails if the queried data cannot be decoded, or if there is an internal MongoDB error.
//
func (repository *Repository) Find(ctx context.Context, f interfaces.Filter) ([]interface{}, error) {
	collection := getCollection(repository)

	cursor, err := collection.Find(ctx, f)

	if err != nil {
		return nil, err
	}

	return decodeCursor(repository, cursor)
}

// FindByID returns a document in the Repository's collection that has a given ID.
//
// It accesses the Repository's collection using the provided ID, and decodes
// the result into the Repository's base type.
//
// It fails if,
//  0. there is an internal MongoDB error, or
//  1. no object is found (in which case ErrorObjectNotFound is returned), or
//  2. multiple objects are found (in which case ErrorMultipleMatches is returned).
//
func (repository *Repository) FindByID(ctx context.Context, id uuid.UUID) (interface{}, error) {
	matches, err := repository.Find(ctx, map[string]interface{}{
		repository.idField: id,
	})

	if err != nil {
		return nil, err
	}

	switch l := len(matches); {
	case l == 0:
		return nil, ErrorObjectNotFound
	case l > 1:
		return nil, ErrorMultipleMatches
	}

	return matches[0], nil
}

// Exists returns whether an object matching a given filter exists.
//
// It fails if the queried data cannot be decoded, or if there is an internal MongoDB error.
//
func (repository *Repository) Exists(ctx context.Context, f interfaces.Filter) (bool, error) {
	matches, err := repository.Find(ctx, f)
	return len(matches) > 0, err
}

// ExistsByID returns whether at least one object having a given ID exists.
//
// It fails if the queried data cannot be decoded, or if there is an internal MongoDB error.
//
func (repository *Repository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	return repository.Exists(ctx, map[string]interface{}{
		repository.idField: id,
	})
}

// Insert inserts a given structure into the Repository's collection.
//
// It decodes the structure into an object that has the same type as the Repository's
// base type, which is subsequently inserted into the collection.
//
// It fails if obj cannot be decoded into the Repository's base type,
// or if there is an internal MongoDB error.
//
func (repository *Repository) Insert(ctx context.Context, obj interface{}) (*mongo.InsertOneResult, error) {
	collection := getCollection(repository)

	newObj, err := decodeIntoBase(repository, obj)

	if err != nil {
		return nil, err
	}

	return collection.InsertOne(ctx, newObj)
}

// InsertMany inserts a given variadic number of structures into the Repository's collection.
//
// It decodes the structures into objects that have the same type as the Repository's
// base type, which are subsequently inserted into the collection.
//
// It fails if an element of obj cannot be decoded into the Repository's base type,
// or if there is an internal MongoDB error.
//
func (repository *Repository) InsertMany(ctx context.Context, obj ...interface{}) (*mongo.InsertManyResult, error) {
	collection := getCollection(repository)

	decoded := make([]interface{}, len(obj))

	for i, o := range obj {
		newObj, err := decodeIntoBase(repository, o)

		if err != nil {
			return nil, err
		}

		decoded[i] = newObj
	}

	return collection.InsertMany(ctx, decoded)
}

// Update updates at most one document in the Repository's collection
// matching a given filter, using a given map of changes.
//
// It decodes the provided map into an object that has the same type as the Repository's base type, which is
// subsequently used to update the first document matching the provided filter.
//
// It fails if there is an internal MongoDB error.
//
func (repository *Repository) Update(ctx context.Context, f interfaces.Filter, changes map[string]interface{}) (*mongo.UpdateResult, error) {
	collection := getCollection(repository)

	filterMap(repository, changes)

	if len(changes) == 0 {
		return &mongo.UpdateResult{
			MatchedCount:  0,
			ModifiedCount: 0,
			UpsertedCount: 0,
			UpsertedID:    nil,
		}, nil
	}

	updates := bson.D{
		{
			"$set", changes,
		},
	}

	return collection.UpdateOne(ctx, f, updates)
}

// UpdateByID updates at most one document in the Repository's collection that has
// a given ID.
//
// It falls back to the Repository's Update method, using the provided ID as a filter.
//
// It fails if there is an internal MongoDB error.
//
func (repository *Repository) UpdateByID(ctx context.Context, id uuid.UUID, changes map[string]interface{}) (*mongo.UpdateResult, error) {
	return repository.Update(ctx, map[string]interface{}{
		repository.idField: id,
	}, changes)
}

// Delete deletes at most one document in the Repository's collection matching a given filter.
//
// It fails if there is an internal MongoDB error.
//
func (repository *Repository) Delete(ctx context.Context, f interfaces.Filter) (*mongo.DeleteResult, error) {
	collection := getCollection(repository)

	return collection.DeleteOne(ctx, f)
}

// DeleteMany deletes all documents in the Repository matching a given filter.
//
// It fails if there is an internal MongoDB error.
//
func (repository *Repository) DeleteMany(ctx context.Context, f interfaces.Filter) (*mongo.DeleteResult, error) {
	collection := getCollection(repository)

	return collection.DeleteMany(ctx, f)
}

// DeleteByID deletes at most one document in the Repository's collection having a given ID.
//
// It fails if there is an internal MongoDB error.
//
func (repository *Repository) DeleteByID(ctx context.Context, id uuid.UUID) (*mongo.DeleteResult, error) {
	return repository.Delete(ctx, map[string]interface{}{
		repository.idField: id,
		repository.idField: id,
	})
}

// getCollection returns a handle to the Repository's collection.
func getCollection(repository *Repository) *mongo.Collection {
	database := repository.Driver.Client.Database(repository.Driver.Database)
	collection := database.Collection(repository.collection)

	return collection
}

// decodeIntoBase takes a map and decodes it into an object that has the same type
// as the Repository's base type.
//
// It fails if the map cannot be decoded.
//
func decodeIntoBase(repository *Repository, obj interface{}) (interface{}, error) {
	t := repository.baseType
	rt := reflect.TypeOf(t).Elem()

	newObj := reflect.New(rt).Interface()

	err := mapstructure.Decode(obj, &newObj)
	return newObj, err
}

// filterMap removes all entries from a given map that are not struct fields of
// the Repository's base type.
func filterMap(repository *Repository, obj map[string]interface{}) {
	t := repository.baseType
	rt := reflect.TypeOf(t).Elem()

	// maps from bson field name to real field name
	fieldMappings := make(map[string]string)

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldName := field.Name

		bsonTag, ok := field.Tag.Lookup("bson")

		switch ok {
		case false:
			fieldMappings[fieldName] = fieldName
		case true:
			bsonName := strings.Split(strings.Trim(bsonTag, " "), ",")[0]
			fieldMappings[bsonName] = fieldName
		}
	}

	for k := range obj {
		realName := fieldMappings[k]

		_, ok := rt.FieldByName(realName)

		if !ok {
			delete(obj, k)
		}
	}
}

// decodeCursor takes a cursor and decodes it into a slice of objects that
// have the same type as the Repository's base type.
//
// It fails if the cursor cannot be fully decoded.
//
func decodeCursor(repository *Repository, cursor *mongo.Cursor) ([]interface{}, error) {
	var matches []interface{}

	t := repository.baseType
	rt := reflect.TypeOf(t).Elem()

	for cursor.Next(context.Background()) {
		r := reflect.New(rt).Interface()

		if err := cursor.Decode(r); err != nil {
			return nil, err
		}

		matches = append(matches, r)
	}

	return matches, nil
}
