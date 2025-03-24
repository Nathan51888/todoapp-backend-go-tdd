package inmemory

import (
	"errors"
	"mytodoapp/domain/todo"
	"slices"

	"github.com/google/uuid"
)

type InMemoryTodoStore struct {
	Todos []todo.Todo
}

func NewInMemoryTodoStore() (*InMemoryTodoStore, error) {
	return &InMemoryTodoStore{}, nil
}

func (i *InMemoryTodoStore) GetTodoAll(userId uuid.UUID) ([]todo.Todo, error) {
	return i.Todos, nil
}

func (i *InMemoryTodoStore) GetTodoByTitle(userId uuid.UUID, title string) (todo.Todo, error) {
	var result todo.Todo
	for _, todo := range i.Todos {
		if todo.Title == title {
			result = todo
		}
	}
	return result, nil
}

func (i *InMemoryTodoStore) GetTodoById(userId uuid.UUID, todoId uuid.UUID) (todo.Todo, error) {
	var result todo.Todo
	for _, item := range i.Todos {
		if item.Id == todoId {
			result = item
			return result, nil
		}
	}

	return todo.Todo{}, errors.New("todo not found")
}

func (i *InMemoryTodoStore) CreateTodo(userId uuid.UUID, title string) (todo.Todo, error) {
	newId := uuid.New()
	createdTodo := todo.Todo{Id: newId, Title: title, Completed: false, UserId: userId}
	i.Todos = append(i.Todos, createdTodo)
	return createdTodo, nil
}

func (i *InMemoryTodoStore) UpdateTodoTitle(userId uuid.UUID, todoId uuid.UUID, title string) (todo.Todo, error) {
	var result todo.Todo
	for index, todo := range i.Todos {
		if todo.Id == todoId {
			i.Todos[index].Title = title
			result = i.Todos[index]
		}
	}
	return result, nil
}

func (i *InMemoryTodoStore) UpdateTodoStatus(userId uuid.UUID, todoId uuid.UUID, completed bool) (todo.Todo, error) {
	var result todo.Todo
	for index, todo := range i.Todos {
		if todo.Id == todoId {
			i.Todos[index].Completed = completed
			result = i.Todos[index]
		}
	}
	return result, nil
}

func (i *InMemoryTodoStore) UpdateTodoById(userId uuid.UUID, todoId uuid.UUID, changedTodo todo.Todo) (todo.Todo, error) {
	var result todo.Todo
	for index, item := range i.Todos {
		if item.Id == todoId {
			i.Todos[index] = changedTodo
			result = i.Todos[index]
		}
	}
	return result, nil
}

func (i *InMemoryTodoStore) DeleteTodoById(userId uuid.UUID, todoId uuid.UUID) (todo.Todo, error) {
	for index, item := range i.Todos {
		if item.Id == todoId {
			i.Todos = slices.Delete(i.Todos, index, index+1)
			return todo.Todo{Id: todoId, Title: item.Title, Completed: false, UserId: userId}, nil
		}
	}
	return todo.Todo{}, errors.New("todo not found")
}
