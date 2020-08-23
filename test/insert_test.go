package test

import (
    "context"
    "fmt"
    "testing"
)

func TestRepository_Insert(t *testing.T) {
    result, err := testRepository.Insert(context.Background(), obiMapKenobi)

    if err != nil {
        t.Error(err)
    }

    fmt.Println(result)
}

func TestRepository_InsertMany(t *testing.T) {
    result, err := testRepository.InsertMany(context.Background(), obiMapKenobi, aniMapSkywalker)

    if err != nil {
        t.Error(err)
    }

    fmt.Println(result)
}
