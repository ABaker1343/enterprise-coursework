package repository

import (
	"testing"
    "encoding/json"
	"github.com/stretchr/testify/assert"
)

type Test struct {
    Input Track
    Output int64
}

func TestUnitAddNewTrack(t *testing.T) {
    // initialize and clear the repo
    Init()
    Clear()
    Create()

    tests := []Test{
        {Input : Track{Id : "Big Band", Audio : "bigbandaudio"}, Output : 1},
        {Input : Track{Id : "wrong band", Audio : "wrongbandaudio"}, Output : 1},
        {Input : Track{Id : "Big Band", Audio : "bigbandaudio"}, Output : 0},
        {Input : Track{Id : "", Audio : "somenewaudio"}, Output : -1},
    }

    for _, test := range tests {
        output := AddNewTrack(test.Input)
        testString, err := json.Marshal(test)
        if err != nil {
            panic("failed to marshal test into json")
        }
        assert.Equal(t, test.Output, output, "new track unit test : " + string(testString))
    }
}

func TestUnitUpdateTrack(t *testing.T) {
    Init()
    Clear()
    Create()

    // insert a new track to the repo
    AddNewTrack(Track{Id : "big band" , Audio : "bigbandaudio"})

    tests := []Test {
        {Input : Track{Id : "big band", Audio : "bigbandaudio"}, Output : 1},
        {Input : Track{Id : "wrong band", Audio : "wrongbandaudio"}, Output : 0},
    }

    for _, test := range tests {
        output := UpdateTrack(test.Input)
        assert.Equal(t, test.Output, output, "update track unit test")
    }
}
