package httpserver_test

import (
	"mytodoapp/adapters/httpserver"
	"mytodoapp/domain/todo"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TodoStoreContract struct {
	NewTodoStore func() httpserver.TodoStore
}

func (c TodoStoreContract) Test(t *testing.T) {
	t.Run("can get todo by title from database", func(t *testing.T) {
		sut := c.NewTodoStore()

		want := todo.Todo{Title: "Todo_test", Completed: "false"}

		got, err := sut.GetTodoByTitle("Todo_test")
		assert.NoError(t, err)
		assert.Equal(t, want, got)

		want = todo.Todo{Title: "Todo_new", Completed: "false"}

		newTodo, err := sut.CreateTodo("Todo_new")
		assert.NoError(t, err)
		assert.Equal(t, want, newTodo)
		got, err = sut.GetTodoByTitle("Todo_new")
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestInMemoryTodoStore(t *testing.T) {
	TodoStoreContract{NewTodoStore: func() httpserver.TodoStore {
		return httpserver.NewInMemoryTodoStore()
	}}.Test(t)
}
