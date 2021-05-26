package mongodialect

import (
    "context"
    "errors"
    "fmt"
    "github.com/daspoet/mongodialect/interfaces"
    "github.com/mitchellh/mapstructure"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "reflect"
    "strings"
)

var (
    // ErrDriverNil occurs when a given driver is nil.
    ErrDriverNil = errors.New("driver must no be nil")

    // ErrCollectionEmpty occurs when the name of a given collection is empty.
    ErrCollectionEmpty = errors.New("collection must not be empty")

    // ErrDocumentNotFound occurs when a lookup does not yield a result.
    ErrDocumentNotFound = errors.New("document was not found")

    // ErrMultipleMatches occurs when a lookup
    // using a given ID yields more than one result.
    ErrMultipleMatches = errors.New("multiple matches for id")
)

// A Repository wraps the Driver and provides
// functionality for performing operations on
// the collections of a MongoDB database. It
// can only access one collection at a time,
// though.
//
// Repository implements the Repository interface.
//
// To allow for generic access to collections
// containing arbitrary data, the Repository
// must be provided with the underlying type
// of the data structure contained in the
// specific collection to access.
type Repository struct {
    baseType   reflect.Type // the type of the data structure stored in the collection; must be a pointer
    idField    string       // the name of the field containing the underlying value's id
    collection string       // the name of the collection to access
    Driver     *Driver      // the Driver used to connect to the database
}

// NewRepository returns a new Repository upon
// validating the given base type and Driver.
//
// If idField is an empty string, the
// default Mongo id ("_id") is used instead.
//
// It fails if the driver is nil, or if
// the provided base type is not a pointer.
//
// It also fails if collection is an empty string.
func NewRepository(baseType reflect.Type, driver *Driver, collection string, idField string) (*Repository, error) {
    if driver == nil {
        return nil, ErrDriverNil
    }

    if collection == "" {
        return nil, ErrCollectionEmpty
    }

    if kind := baseType.Kind(); kind != reflect.Ptr {
        return nil, fmt.Errorf("baseType must be a pointer, not '%s'", kind)
    }

    // fallback to Mongo's default id
    if idField == "" {
        idField = "_id"
    }

    return &Repository{
        baseType:   baseType,
        idField:    idField,
        collection: collection,
        Driver:     driver,
    }, nil
}

// InitialiseNewRepository builds all components
// needed for and combines them into a Repository.
//
// see NewRepository
func InitialiseNewRepository(baseType reflect.Type, port uint, hostname, database, collection, idField string) (*Repository, error) {
    url := NewDatabaseURL(hostname, port)
    driver := NewDriver(url, database)

    if err := driver.OpenConnection(context.Background()); err != nil {
        return nil, err
    }
    return NewRepository(baseType, driver, collection, idField)
}

// Type returns a pointer to r's base type, which has
// the same type as the data stored in r's collection.
func (r *Repository) Type() reflect.Type {
    return r.baseType
}

// Find finds all documents in r's collection matching
// f and decodes each into a value of r's base type.
//
// It fails if the queried data cannot be decoded,
// or if there is an internal MongoDB error.
func (r *Repository) Find(ctx context.Context, f interfaces.Filter) ([]interface{}, error) {
    cursor, err := collection(r).Find(ctx, f)
    if err != nil {
        return nil, err
    }
    return decodeCursor(r, cursor)
}

// FindByID finds a document in r's collection that has the
// given id and decodes it into a value of r's base type.
//
// It fails if
//
//  1. there is an internal MongoDB error (in which
//     case the respective error is returned), or
//
//  2. if no document is found (in which case
//     ErrDocumentNotFound is returned), or
//
//  3. if multiple documents are found (in which
//     case ErrMultipleMatches is returned).
//
func (r *Repository) FindByID(ctx context.Context, id interface{}) (interface{}, error) {
    matches, err := r.Find(ctx, map[string]interface{}{
        r.idField: id,
    })

    if err != nil {
        return nil, err
    }

    switch l := len(matches); {
    case l == 0:
        return nil, ErrDocumentNotFound
    case l > 1:
        return nil, ErrMultipleMatches
    }
    return matches[0], nil
}

// Exists returns whether a document
// matching f exists in r's collection.
//
// It fails if the queried data cannot be decoded,
// or if there is an internal MongoDB error.
func (r *Repository) Exists(ctx context.Context, f interfaces.Filter) (bool, error) {
    matches, err := r.Find(ctx, f)
    return len(matches) > 0, err
}

// ExistsByID returns whether at least one document
// having the given id exists in r's collection.
//
// It fails if the queried data cannot be decoded,
// or if there is an internal MongoDB error.
func (r *Repository) ExistsByID(ctx context.Context, id interface{}) (bool, error) {
    return r.Exists(ctx, map[string]interface{}{
        r.idField: id,
    })
}

// Insert inserts a value into r's collection.
//
// It decodes v into a value of r's base type,
// which is subsequently inserted into r's collection.
//
// It fails if v cannot be decoded into r's base
// type, or if there is an internal MongoDB error.
func (r *Repository) Insert(ctx context.Context, v interface{}) (*mongo.InsertOneResult, error) {
    dec, err := decodeIntoBase(r, v)
    if err != nil {
        return nil, err
    }
    return collection(r).InsertOne(ctx, dec)
}

// InsertMany inserts a variadic number
// of values into r's collection.
//
// It decodes each element in v into a value of
// r's base type, which is subsequently inserted
// into r's collection.
//
// It fails if an element of v cannot be decoded into
// r's base type, or if there is an internal MongoDB error.
func (r *Repository) InsertMany(ctx context.Context, v ...interface{}) (*mongo.InsertManyResult, error) {
    decoded := make([]interface{}, len(v))
    for i, o := range v {
        dec, err := decodeIntoBase(r, o)
        if err != nil {
            return nil, err
        }
        decoded[i] = dec
    }
    return collection(r).InsertMany(ctx, decoded)
}

// Update updates at most one document in r's
// collection matching f, using the given changes.
//
// It decodes changes into a value of r's
// base type, which is subsequently used
// to update the first document matching f.
//
// It fails if there is an internal MongoDB error.
func (r *Repository) Update(ctx context.Context, f interfaces.Filter, changes map[string]interface{}) (*mongo.UpdateResult, error) {
    filterMap(r, changes)
    if len(changes) == 0 {
        return &mongo.UpdateResult{
            MatchedCount:  0,
            ModifiedCount: 0,
            UpsertedCount: 0,
            UpsertedID:    nil,
        }, nil
    }
    updates := bson.D{{"$set", changes}}
    return collection(r).UpdateOne(ctx, f, updates)
}

// UpdateByID updates at most one document
// in r's collection that has the given id.
//
// It falls back to the Repository's Update
// method, using the provided id as a filter.
//
// It fails if there is an internal MongoDB error.
func (r *Repository) UpdateByID(ctx context.Context, id interface{}, changes map[string]interface{}) (*mongo.UpdateResult, error) {
    return r.Update(ctx, map[string]interface{}{
        r.idField: id,
    }, changes)
}

// Delete deletes at most one document
// in r's collection matching f.
//
// It fails if there is an internal MongoDB error.
func (r *Repository) Delete(ctx context.Context, f interfaces.Filter) (*mongo.DeleteResult, error) {
    return collection(r).DeleteOne(ctx, f)
}

// DeleteMany deletes all documents
// in r's collection matching f.
//
// It fails if there is an internal MongoDB error.
func (r *Repository) DeleteMany(ctx context.Context, f interfaces.Filter) (*mongo.DeleteResult, error) {
    return collection(r).DeleteMany(ctx, f)
}

// DeleteByID deletes at most one document
// in r's collection having the given id.
//
// It fails if there is an internal MongoDB error.
func (r *Repository) DeleteByID(ctx context.Context, id interface{}) (*mongo.DeleteResult, error) {
    return r.Delete(ctx, map[string]interface{}{
        r.idField: id,
    })
}

// collection returns a handle for r's collection.
func collection(r *Repository) *mongo.Collection {
    db := r.Driver.Client.Database(r.Driver.Database)
    return db.Collection(r.collection)
}

// decodeIntoBase decodes v into a new value of r's
// base type. v must be a pointer to a map or struct.
func decodeIntoBase(r *Repository, v interface{}) (interface{}, error) {
    el := r.baseType.Elem()
    dec := reflect.New(el).Interface()

    err := mapstructure.Decode(v, &dec)
    return dec, err
}

// filterMap removes all entries from v that
// do not refer to fields of r's base type.
func filterMap(r *Repository, v map[string]interface{}) {
    el := r.baseType.Elem()

    // maps from bson field name to real field name
    fieldMappings := make(map[string]string)

    for i := 0; i < el.NumField(); i++ {
        field := el.Field(i)

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
    for k := range v {
        realName := fieldMappings[k]
        _, ok := el.FieldByName(realName)
        if !ok {
            delete(v, k)
        }
    }
}

// decodeCursor decodes cur into a
// slice of values of r's base type.
//
// It fails if cur cannot be fully decoded.
func decodeCursor(r *Repository, cur *mongo.Cursor) ([]interface{}, error) {
    var matches []interface{}
    el := r.baseType.Elem()

    for cur.Next(context.Background()) {
        r := reflect.New(el).Interface()
        if err := cur.Decode(r); err != nil {
            return nil, err
        }
        matches = append(matches, r)
    }
    return matches, nil
}
