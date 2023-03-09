package resources

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "tracks/repository"
)

func setUpRepo() {
    repository.Init()
    repository.Clear()
    repository.Create()
}

func TestAddNewTrack(t *testing.T) {
    setUpRepo()
}
