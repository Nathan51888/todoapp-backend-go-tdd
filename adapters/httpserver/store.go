package httpserver

import "mytodoapp/domain/todo"

type TodoStore interface {
	GetTodoByTitle(title string) (todo.Todo, error)
	CreateTodo(title string) (todo.Todo, error)
}
