package todo

import "github.com/google/uuid"

type TodoStore interface {
	GetTodoAll() ([]Todo, error)
	GetTodoByTitle(title string) (Todo, error)
	GetTodoById(todoId uuid.UUID) (Todo, error)
	CreateTodo(title string) (Todo, error)
	UpdateTodoTitle(todoId uuid.UUID, title string) (Todo, error)
	UpdateTodoStatus(todoId uuid.UUID, completed bool) (Todo, error)
	UpdateTodoById(todoId uuid.UUID, changedTodo Todo) (Todo, error)
	DeleteTodoById(todoId uuid.UUID) (Todo, error)
}
