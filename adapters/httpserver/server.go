package httpserver

import (
	"mytodoapp/adapters/httpserver/handler"
	"mytodoapp/adapters/httpserver/middleware"
	"mytodoapp/domain/todo"
	"mytodoapp/domain/user"
	"net/http"
)

type TodoServer struct {
	http.Handler
}

func NewTodoServer(todoStore todo.TodoStore, userStore user.UserStore) *TodoServer {
	server := new(TodoServer)

	mux := http.NewServeMux()
	handler.NewTodoHandler(mux, todoStore, userStore)
	handler.NewUserHandler(mux, userStore)
	handler.NewAuthHandler(mux, userStore)
	stack := middleware.CreateStack(
		middleware.AllowCors,
		middleware.RecoveryMiddleware,
	)
	handler := stack(mux)

	server.Handler = handler

	return server
}
