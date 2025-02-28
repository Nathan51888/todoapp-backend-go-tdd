package main

import (
	"log"
	"mytodoapp/adapters/httpserver"
	"mytodoapp/adapters/persistence/inmemory"
	"mytodoapp/domain/todo"
	"net/http"
)

func main() {
	store := &inmemory.InMemoryTodoStore{Todos: []todo.Todo{
		{Title: "Todo1", Completed: "false"},
	}}
	server := httpserver.NewTodoServer(store)
	log.Fatal(http.ListenAndServe(":8080", server))
}
