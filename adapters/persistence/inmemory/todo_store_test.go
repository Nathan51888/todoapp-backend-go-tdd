package inmemory_test

import (
	"mytodoapp/adapters/persistence/inmemory"
	"mytodoapp/domain/todo"
	"testing"
)

func TestInMemoryTodoStore(t *testing.T) {
	todo.TodoStoreContract{NewTodoStore: func() todo.TodoStore {
		return inmemory.NewInMemoryTodoStore()
	}}.Test(t)
}
