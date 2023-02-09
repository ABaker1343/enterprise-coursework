package repository

import (

)

var repo map[string]interface{}

func Init() {
    repo = map[string]interface{} {}
}

func AddNewTrack(id string, dataString string) int {
    // function that adds a track to the repository
    if _, ok := repo[id]; ok {
        // track already exists in repository
        return 0
    }

    repo[id] = dataString
    return 1
}

func GetTrackById(id string) (string, int) {
    // function that gets all the tracks in the repo and returns the audio base64 encoded
    
    data, ok := repo[id]
    if !ok {
        //track does not exist
        return "", 0
    }

    return data.(string), 1
}

func GetAllTracks() ([]string, int) {
    // returns all the tracks in the repo
    trackList := make([]string, 0)
    for k := range repo {
        trackList = append(trackList, k)
    }
    return trackList, len(trackList)
}
