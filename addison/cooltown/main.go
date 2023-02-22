package main

import (
	"net/http"
    "cooltown/resources"
    "fmt"
)

func main() {
    fmt.Println("this is the cooltown microservice")

    resources.Init()

    fmt.Println("listening on port 3002")
    http.ListenAndServe(":3002", resources.Router())

}
