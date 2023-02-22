package resources

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
    "strings"
	"github.com/gorilla/mux"
)

var logger *log.Logger

func Init() {
    logger = log.New(os.Stdin, "", 0)
}

func Router() http.Handler {
    router := mux.NewRouter()
    
    router.HandleFunc("/cooltown", fetchTrack).Methods("POST")

    return router
}

func fetchTrack(w http.ResponseWriter, r *http.Request) {
    // function that takes a POST request and uses the other services to return the audio
    // for the song snippet

    // get the audio
    t := map[string]interface{} {}
    if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
        // failed to decode body
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    audio, ok := t["Audio"]
    if !ok || audio == "" {
        // Audio is not a field in the json
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    searchRequestBody, err := json.Marshal(
        map[string]interface{} {"Audio" : audio},
    )
    if err != nil {
        // failed to marshal response
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    // create request to find the song name from the search service on port 3001
    searchRes, err := http.Post("http://127.0.0.1:3001/search", "application/json", bytes.NewBuffer(searchRequestBody))
    if err != nil {
        // failed http request to search microservice
        w.WriteHeader(http.StatusInternalServerError)
        logger.Output(2, "failed http request to search microservice: " + err.Error())
        return
    }

    if searchRes.StatusCode != http.StatusOK {
        w.WriteHeader(searchRes.StatusCode)
        logger.Output(2, "response code from searching service: " + searchRes.Status)
        return
    }

    // take the name of the song and fetch the audio from the tracks api
    defer searchRes.Body.Close()

    searchBody := map[string]interface{} {}
    err = json.NewDecoder(searchRes.Body).Decode(&searchBody)
    if err != nil {
        // failed to decode body of search response
        logger.Output(2, "failed to decode response from search microservice: " + err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // make request to tracks
    trackId, ok := searchBody["Id"]
    if !ok {
        // search response did not contain Id
        w.WriteHeader(http.StatusInternalServerError)
        logger.Output(2, "search response did not contain Id")
        return
    }
    trackURL := "http://127.0.0.1:3000/tracks/" + strings.Replace(trackId.(string), " ", "+", -1)
    //trackURL := "nice" + trackId.(string)
    tracksRes, err := http.Get(trackURL)
    
    if tracksRes.StatusCode != http.StatusOK {
        w.WriteHeader(tracksRes.StatusCode)
        logger.Output(2, "response from tracks api: " + tracksRes.Status)
        return
    }

    defer tracksRes.Body.Close()

    tracksBody := map[string]interface{} {}
    err = json.NewDecoder(tracksRes.Body).Decode(&tracksBody)
    if err != nil {
        // failed to decode the tracks body
        logger.Output(2, "failed to decode tracks body: " + err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    trackAudio, ok := tracksBody["Audio"]
    if !ok {
        // tracks did not respond with audio
        w.WriteHeader(http.StatusInternalServerError)
        logger.Output(2, "tracks did not respond with audio")
        return
    }

    response := map[string]interface{} {"Audio" : trackAudio}
    err = json.NewEncoder(w).Encode(response)
    if err != nil {
        // failed to encode response
        logger.Output(2, "failed to encode response")
        w.WriteHeader(http.StatusOK)
        return
    }
    w.WriteHeader(http.StatusOK)
    
}
