// tracks micro service

// this micorservice should have the following 3 use cases:
// put tracks onto the list
// retrieve tracks out of the list by name
// list all the tracks in the list

// JSON format for this microservices is
// {"Id": <int>, "Audio": <base64 string>}

package main

import (
	"net/http"
	"tracks/repository"
	"tracks/resources"
)

func main() {
    print("this is the tracks microservice\n")
    repository.Init()

    http.ListenAndServe(":3000", resources.Router())
}
