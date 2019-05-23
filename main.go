package main

import (
	"log"
	"net/http"
)

func main() {

	router := newRouter()

	server := http.ListenAndServe(":8082", router)

	log.Fatal(server)

}
