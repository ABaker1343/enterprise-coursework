package resources

import (
    "net/http"
    "tracks/repository"
    "github.com/gorilla/mux"
    "encoding/json"
    "fmt"
)

func Router() http.Handler {
    router := mux.NewRouter()

    // store new tracks
    router.HandleFunc("/tracks/{id}", addNewTrack).Methods("PUT")

    // retrieve tracks by ID
    router.HandleFunc("/tracks/{id}", getTrackById).Methods("GET")

    // list all available tracks
    router.HandleFunc("/tracks", getAllTracks).Methods("GET")
    
    return router
}

func addNewTrack(w http.ResponseWriter, r *http.Request) {
    // function that takes a http request and adds a new track to the track list
    data := map[string]interface{} {}

    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    }

    id, ok := data["Id"]
    if !ok {
        // missing Id field
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    audio, ok := data["Audio"]
    if !ok {
        // missing audio field
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    if audio == "" {
        //audio is empty
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    track := repository.Track{Id : id.(string), Audio : audio.(string)}

    repoResponse := repository.AddNewTrack(track)
    
    if repoResponse == 0 {
        // track already exists
        w.WriteHeader(http.StatusConflict)
        return
    } else if repoResponse == -1 {
        // unexpected error
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    response := map[string]interface{} {"Id" : id, "Audio" : audio}
    if err := json.NewEncoder(w).Encode(response); err != nil {
        // failed to encode data
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)

}

func getTrackById(w http.ResponseWriter, r *http.Request) {
    //function to get a single track by its id (name)

    data := mux.Vars(r)
    

    id, ok := data["id"]
    if !ok {
        w.WriteHeader(http.StatusBadRequest)
    }

    trackData, repoResponse := repository.GetTrackById(id)
    if repoResponse == 0 {
        //track does not exist
        w.WriteHeader(http.StatusNotFound)
        return
    }
    if repoResponse == -1 {
        // something unexpected went wrong
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    response := map[string]interface{} {"Id": id, "Audio": trackData}
    if err := json.NewEncoder(w).Encode(response); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Println(err)
    }
    w.WriteHeader(http.StatusOK)

}

func getAllTracks(w http.ResponseWriter, r *http.Request) {
    // function that will retrieve all tracks in the repo

    allTracks, numTracks := repository.GetAllTracks()

    if numTracks == 0 {
        // no tracks in list
        w.WriteHeader(http.StatusNotFound)
    }

    titles := make([]string, len(allTracks))
    for _, t := range allTracks {
        titles = append(titles, t.Id)
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(titles)
}
