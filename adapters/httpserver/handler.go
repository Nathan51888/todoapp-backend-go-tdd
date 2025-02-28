package httpserver

import (
	"fmt"
	"mytodoapp/domain/interactions"
	"net/http"
)

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/test", handleTodo)
	return mux
}

func handleTodo(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Fprint(w, interactions.Greet(name))
}
