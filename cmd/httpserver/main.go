package main

import (
	"log"
	"mytodoapp/adapters/httpserver"
	"mytodoapp/adapters/persistence/postgre"
	"net/http"
)

func main() {
	store := postgre.NewPostgreTodoStore()
	server := httpserver.NewTodoServer(store)
	log.Fatal(http.ListenAndServe(":8080", server))
}
