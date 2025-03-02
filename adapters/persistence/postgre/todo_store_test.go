package postgre_test

import (
	"mytodoapp/adapters/persistence/postgre"
	"mytodoapp/domain/todo"
	"testing"
)

func TestPostgreTodoStore(t *testing.T) {
	todo.TodoStoreContract{NewTodoStore: func() todo.TodoStore {
		return postgre.NewPostgreTodoStore("postgres://postgres:test@localhost:5433/postgres")
	}}.Test(t)
}
