package httpserver

import (
	"encoding/json"
	"fmt"
	"log"
	"mytodoapp/adapters/persistence/inmemory"
	"mytodoapp/domain/todo"
	"net/http"
)

type TodoService interface {
	GetTodoByTitle(title string) (todo.Todo, error)
}

type TodoServer struct {
	store todo.TodoStore
	http.Handler
}

func NewTodoServer(store todo.TodoStore) *TodoServer {
	server := new(TodoServer)

	server.store = store

	mux := http.NewServeMux()
	mux.HandleFunc("/todo", server.GetTodoByTitle)

	server.Handler = mux

	return server
}

func (t *TodoServer) GetTodoByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	result, err := t.store.GetTodoByTitle(title)
	if err != nil {
		log.Fatalln(err)
	}

	json.NewEncoder(w).Encode(&result)
}

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/todo", handleTodo)
	return mux
}

func handleTodo(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	result, err := inmemory.NewInMemoryTodoStore().GetTodoByTitle(title)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprint(w, todo.GetTodoTitle(result))
}
