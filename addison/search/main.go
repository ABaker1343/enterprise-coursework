package main

import(
    "fmt"
    "net/http"
    "search/resources"
    "log"
)

func main() {
    fmt.Println("this is the searching microservice")
    
    if resources.Init() != 0 {
        log.Fatal("failed to initialize resources")
    }
    fmt.Println("listening on port 3001")
    http.ListenAndServe(":3001", resources.Router())

}
