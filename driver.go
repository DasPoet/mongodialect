package mongodialect

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/url"
	"strings"
)

// A Driver is a wrapper for connecting to a MongoDB database.
type Driver struct {
	URL      url.URL                // the URL
	Database string                 // the name of the database
	Client   *mongo.Client          // a handle representing a pool of connections to the database
	Options  *options.ClientOptions // options for configuring the client
}

// NewDriver returns a new driver.
func NewDriver(url url.URL, database string) *Driver {
	return &Driver{
		URL:      url,
		Database: database,
		Options:  options.Client(),
	}
}

// OpenConnection establishes a connection to the database.
//
// It fails if there is an internal MongoDB error.
//
func (driver *Driver) OpenConnection(ctx context.Context) error {
	var err error

	opts := driver.Options

	if strings.EqualFold(opts.GetURI(), "") {
		uri := driver.URL.String()
		opts = opts.ApplyURI(uri)
	}

	driver.Client, err = mongo.Connect(ctx, opts)

	return err
}

// CloseConnection disconnects from the database.
//
// It fails if there is an internal MongoDB error.
//
func (driver *Driver) CloseConnection(ctx context.Context) error {
	if driver.Client == nil {
		return nil
	}

	return driver.Client.Disconnect(ctx)
}

// IsAlive returns whether the connection to the database is still alive.
func (driver *Driver) IsAlive(ctx context.Context) bool {
	if driver.Client == nil {
		return false
	}

	return driver.Client.Ping(ctx, nil) != nil
}
