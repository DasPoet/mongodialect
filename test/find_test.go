package test

import (
    "context"
    "fmt"
    "github.com/google/uuid"
    "testing"
)

func TestRepository_Exists(t *testing.T) {
    ok, err := testRepository.Exists(context.Background(), map[string]interface{}{
        "surname": "Kenobi",
    })

    if err != nil {
        t.Error(err)
    }

    fmt.Println(ok)
}

func TestRepository_ExistsByID(t *testing.T) {
    ok, err := testRepository.ExistsByID(context.Background(), uuid.MustParse("822989de-fb98-40f8-a7e3-459277108b67"))

    if err != nil {
        panic(err)
    }

    fmt.Println(ok)
}

func TestRepository_Find(t *testing.T) {
    result, err := testRepository.Find(context.Background(), map[string]interface{}{
        "name": "General",
    })

    if err != nil {
        t.Error(err)
    }

    for _, i := range result {
        fmt.Println(i)
    }
}

func TestRepository_FindByID(t *testing.T) {
    result, err := testRepository.FindByID(context.Background(), uuid.MustParse("f2c17a29-6ff6-4c9a-914f-ae635c84db89"))

    if err != nil {
        t.Error(err)
    }

    fmt.Println(result)
}
