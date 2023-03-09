package resources

import (
	"encoding/json"
	"net/http"
	"os"
	"bytes"
	"io"
	"log"

	//"strings"  -- used to read from file which has been removed
	"github.com/gorilla/mux"
)

type APIResponse struct {
    Status string `json:"status"`
    Result APIResult `json:"result"`
}

type APIResult struct {
    Artist string `json:"artist"`
    Title string `json:"title"`
}

var token string
const KEY = "test"
var logger *log.Logger

func Init() int {
    // this function is for initialisation and will read the api key from the api_token.txt file
    // for the coursework the token will be hard coded in by the marker and so this function can be ignored

    // read the token from a file
    /*if data, err := os.ReadFile("api_token.txt"); err == nil {
        token = string(data)
        token = strings.TrimSuffix(token, "\n")
        logger = log.New(os.Stdin, "", 0)
    } else {
        // failed to read token
        log.Fatal("failed to read api token")
    }*/

    token = KEY
    logger = log.New(os.Stdin, "", 0)
    return 0
}

func Router() http.Handler {
    router := mux.NewRouter()

    // searching url
    router.HandleFunc("/search", search).Methods("POST")

    return router
}

func search(w http.ResponseWriter, r *http.Request) {
    // function that will take http request and search audd.io

    requestBody := map[string]interface{} {}

    if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
        // failed to decode json
        w.WriteHeader(http.StatusInternalServerError)
        logger.Output(2, "failed to decode request body: " + err.Error())
    }

    audio, ok := requestBody["Audio"]
    if !ok {
        // audio is missing
        w.WriteHeader(http.StatusBadRequest)
    }

    // create a request for the audd.io api

    apiReqBody := map[string]interface{} {"api_token" : token, "audio" : audio}
    marshalledBody, err := json.Marshal(apiReqBody)
    if err != nil {
        // failed to marshall body
        logger.Output(2, "failed to marshall api request body: " + err.Error())
        w.WriteHeader(http.StatusInternalServerError)
    }

    apiSendData := bytes.NewBuffer(marshalledBody)
    apiRes, err := http.Post("https://api.audd.io/recognize", "application/json", apiSendData)
    if err != nil {
        // failed http request
        w.WriteHeader(http.StatusInternalServerError)
        logger.Output(2, "http request to audd.io failed: " + err.Error())
    }

    if apiRes.StatusCode != http.StatusOK {
        w.WriteHeader(http.StatusInternalServerError)
        logger.Output(2, "api request failed with code " + apiRes.Status)
    }

    defer apiRes.Body.Close()

    apiResBodyMarshalled, err := io.ReadAll(apiRes.Body)
    if err != nil {
        // failed to read response body
        w.WriteHeader(http.StatusInternalServerError)
        logger.Output(2, "failed to read http response body of api request to audd.io: " + err.Error())
    }

    apiResBody := APIResponse{}
    err = json.Unmarshal(apiResBodyMarshalled, &apiResBody)
    if err != nil {
        // failed to unmarshall body
        w.WriteHeader(http.StatusInternalServerError)
        logger.Output(2, "failed to unmarshall response from audd.io api response: " + err.Error())
    }

    // check the api response for success
    if apiResBody.Status != "success" {
        w.WriteHeader(http.StatusInternalServerError)
        logger.Output(2, "api response error: " + apiResBody.Status)
        logger.Output(2, string(apiResBodyMarshalled))
    }

    // return the track title to the user

    userRes := map[string]interface{} {"Id" : apiResBody.Result.Title}
    err = json.NewEncoder(w).Encode(userRes)
    if err != nil {
        // failed to encode response
        w.WriteHeader(http.StatusInternalServerError)
        logger.Output(2, "failed to Encode response to user: " + err.Error())
    }

    w.WriteHeader(http.StatusOK)

}
