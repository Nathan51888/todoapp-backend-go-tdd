package main

import (
	"log"
	"mytodoapp/adapters/httpserver"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", httpserver.NewHandler()))
}
