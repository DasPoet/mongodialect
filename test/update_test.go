package test

import (
    "context"
    "fmt"
    "github.com/google/uuid"
    "testing"
)

func TestRepository_Update(t *testing.T) {
    result, err := testRepository.Update(context.Background(), map[string]interface{}{
        "surname": "Skywalker",
    }, map[string]interface{}{
        "surname": "Highwalker",
    })

    if err != nil {
        t.Error(err)
    }

    fmt.Println(result)
}

func TestRepository_UpdateByID(t *testing.T) {
    result, err := testRepository.UpdateByID(context.Background(), uuid.MustParse("f2c17a29-6ff6-4c9a-914f-ae635c84db89"), map[string]interface{}{
        "surname": "Skywalker",
    })

    if err != nil {
        t.Error(err)
    }

    fmt.Println(result)
}
