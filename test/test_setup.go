package test

import (
	"context"
	"github.com/google/uuid"
	"github.com/daspoet/mongodialect"
	"github.com/daspoet/mongodialect/interfaces"
	"reflect"
)

const (
	hostname = "localhost"
	database = "testBase"
	port     = 27017
)

type testDocument struct {
	ID      uuid.UUID `bson:"ID, omitempty"`
	Name    string    `bson:"name, omitempty"`
	Surname string    `bson:"surname, omitempty"`
	Age     uint      `bson:"age, omitempty"`
}

var obi = testDocument{
	ID:      uuid.MustParse("822989de-fb98-40f8-a7e3-459277108b67"),
	Name:    "General",
	Surname: "Kenobi",
	Age:     1,
}

var ani = testDocument{
	ID:      uuid.MustParse("f2c17a29-6ff6-4c9a-914f-ae635c84db89"),
	Name:    "General",
	Surname: "Skywalker",
	Age:     0,
}

func makeRepository() interfaces.Repository {
	uri := mongodialect.NewDatabaseURL(hostname, port)
	driver := mongodialect.NewDriver(uri, database)

	if err := driver.OpenConnection(context.Background()); err != nil {
		panic(err)
	}

	base := reflect.TypeOf(new(testDocument))
	repository, err := mongodialect.NewRepository(base, driver, "test", "ID")

	if err != nil {
		panic(err)
	}

	return repository
}

var testRepository = makeRepository()
