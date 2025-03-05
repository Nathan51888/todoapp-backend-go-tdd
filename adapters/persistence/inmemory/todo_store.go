package inmemory

import (
	"mytodoapp/domain/todo"
)

type InMemoryTodoStore struct {
	Todos []todo.Todo
}

func NewInMemoryTodoStore() (*InMemoryTodoStore, error) {
	return &InMemoryTodoStore{}, nil
}

func (i *InMemoryTodoStore) GetTodoAll() ([]todo.Todo, error) {
	return i.Todos, nil
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

func (i *InMemoryTodoStore) UpdateTodoTitle(todoToChange string, title string) (todo.Todo, error) {
	var result todo.Todo
	for index, todo := range i.Todos {
		if todo.Title == todoToChange {
			i.Todos[index].Title = title
			result = i.Todos[index]
		}
	}
	return result, nil
}

func (i *InMemoryTodoStore) UpdateTodoStatus(todoToChange string, completed string) (todo.Todo, error) {
	var result todo.Todo
	for index, todo := range i.Todos {
		if todo.Title == todoToChange {
			i.Todos[index].Completed = completed
			result = i.Todos[index]
		}
	}
	return result, nil
}
