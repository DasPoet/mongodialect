package test

import (
    "context"
    "fmt"
    "github.com/google/uuid"
    "testing"
)

func TestRepository_Delete(t *testing.T) {
    result, err := testRepository.Delete(context.Background(), map[string]interface{}{
        "surname": "Kenobi",
    })

    if err != nil {
        t.Error(err)
    }

    fmt.Println(result)
}

func TestRepository_DeleteByID(t *testing.T) {
    result, err := testRepository.DeleteByID(context.Background(), uuid.MustParse("f2c17a29-6ff6-4c9a-914f-ae635c84db89"))

    if err != nil {
        t.Error(err)
    }

    fmt.Println(result)
}

func TestRepository_DeleteMany(t *testing.T) {
    result, err := testRepository.DeleteMany(context.Background(), map[string]interface{}{
        "name": "General",
    })

    if err != nil {
        t.Error(err)
    }

    fmt.Println(result)
}
