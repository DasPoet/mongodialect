package test

import (
    "context"
    "github.com/fatih/structs"
    "github.com/google/uuid"
    "github.com/informatik-q2/mongodialect"
    "github.com/informatik-q2/mongodialect/interfaces"
)

const (
    hostname = "localhost"
    database = "testBase"
    port     = 27017
)

type testvect struct {
    ID      uuid.UUID `bson:"ID, omitempty"`
    Name    string    `bson:"name, omitempty"`
    Surname string    `bson:"surname, omitempty"`
    Age     uint      `bson:"age, omitempty"`
}

var obi = testvect{
    ID:      uuid.MustParse("822989de-fb98-40f8-a7e3-459277108b67"),
    Name:    "General",
    Surname: "Kenobi",
    Age:     1,
}

var ani = testvect{
    ID:      uuid.MustParse("f2c17a29-6ff6-4c9a-914f-ae635c84db89"),
    Name:    "General",
    Surname: "Skywalker",
    Age:     0,
}

var obiMapKenobi = structs.Map(obi)

var aniMapSkywalker = structs.Map(ani)

func makeRepository() interfaces.Repository {
    uri := mongodialect.NewDatabaseURL(hostname, port)
    driver := mongodialect.NewDriver(uri, database)

    if err := driver.OpenConnection(context.Background()); err != nil {
        panic(err)
    }

    repository, err := mongodialect.NewRepository(&testvect{}, driver, "test", "ID")

    if err != nil {
        panic(err)
    }

    return repository
}

var testRepository = makeRepository()
