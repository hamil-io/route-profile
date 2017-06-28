package main

import (
	"log"
	"net/http"
	"route-profile/api"
)

func main() {

	router := api.MainRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
