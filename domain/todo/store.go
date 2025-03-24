package todo

import "github.com/google/uuid"

type TodoStore interface {
	GetTodoAll(userId uuid.UUID) ([]Todo, error)
	GetTodoByTitle(userId uuid.UUID, title string) (Todo, error)
	GetTodoById(userId uuid.UUID, todoId uuid.UUID) (Todo, error)
	CreateTodo(userId uuid.UUID, title string) (Todo, error)
	UpdateTodoTitle(userId uuid.UUID, todoId uuid.UUID, title string) (Todo, error)
	UpdateTodoStatus(userId uuid.UUID, todoId uuid.UUID, completed bool) (Todo, error)
	UpdateTodoById(userId uuid.UUID, todoId uuid.UUID, changedTodo Todo) (Todo, error)
	DeleteTodoById(userId uuid.UUID, todoId uuid.UUID) (Todo, error)
}
