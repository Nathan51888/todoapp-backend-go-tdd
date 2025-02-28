package httpserver

import "mytodoapp/domain/todo"

type InMemoryTodoStore struct {
	Todos []todo.Todo
}

func NewInMemoryTodoStore() *InMemoryTodoStore {
	return &InMemoryTodoStore{Todos: []todo.Todo{{Title: "Todo_test", Completed: "false"}}}
}

func (i *InMemoryTodoStore) GetTodoByTitle(title string) (todo.Todo, error) {
	var result todo.Todo
	for _, todo := range i.Todos {
		if todo.Title == title {
			result = todo
		}
	}
	return result, nil
}

func (i *InMemoryTodoStore) CreateTodo(title string) (todo.Todo, error) {
	createdTodo := todo.Todo{Title: title, Completed: "false"}
	i.Todos = append(i.Todos, createdTodo)
	return createdTodo, nil
}
