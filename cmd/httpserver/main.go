package main

import (
	"log"
	go_specs_greet "mytodoapp"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(go_specs_greet.Handler)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
