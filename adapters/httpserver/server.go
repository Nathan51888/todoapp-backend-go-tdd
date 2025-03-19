package httpserver

import (
	"mytodoapp/adapters/httpserver/handler"
	"mytodoapp/adapters/httpserver/middleware"
	"mytodoapp/adapters/persistence/inmemory"
	"mytodoapp/domain/todo"
	"net/http"
)

type TodoServer struct {
	http.Handler
}

func NewTodoServer(todoStore todo.TodoStore) *TodoServer {
	server := new(TodoServer)

	mux := http.NewServeMux()
	handler.NewTodoHandler(mux, todoStore, &inmemory.InMemoryUserStore{})
	handler.NewUserHandler(mux, &inmemory.InMemoryUserStore{})
	stack := middleware.CreateStack(
		middleware.AllowCors,
	)
	handler := stack(mux)

	server.Handler = handler

	return server
}
