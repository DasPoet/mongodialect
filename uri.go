package mongodialect

import (
    "fmt"
    "net/url"
)

// NewDatabaseURL creates a new URL for connecting to a MongoDB database.
func NewDatabaseURL(hostname string, port uint) url.URL {
    return url.URL{
        Scheme: "mongodb",
        Host:   fmt.Sprintf("%s:%v", hostname, port),
    }
}
