package httpserver

import (
	"fmt"
	"mytodoapp/domain/todo"
	"net/http"
)

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/todo", handleTodo)
	return mux
}

func handleTodo(w http.ResponseWriter, r *http.Request) {
	// title := r.URL.Query().Get("title")
	result := todo.Todo{Title: "Todo1", Completed: "true"}
	fmt.Fprint(w, todo.GetTodoTitle(result))
}
