package postgre

import "mytodoapp/domain/todo"

type PostgreTodoStore struct{}

func NewPostgreTodoStore() *PostgreTodoStore {
	return &PostgreTodoStore{}
}

func (p *PostgreTodoStore) GetTodoByTitle(title string) (todo.Todo, error) {
	return todo.Todo{}, nil
}

func (p *PostgreTodoStore) CreateTodo(title string) (todo.Todo, error) {
	return todo.Todo{}, nil
}
