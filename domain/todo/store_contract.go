package todo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TodoStoreContract struct {
	NewTodoStore func() TodoStore
}

func (c TodoStoreContract) Test(t *testing.T) {
	t.Run("can create, get, update todo by title from database", func(t *testing.T) {
		sut := c.NewTodoStore()

		want := Todo{Title: "Todo_new", Completed: "false"}
		newTodo, err := sut.CreateTodo("Todo_new")
		assert.NoError(t, err)
		assert.Equal(t, want, newTodo)
		got, err := sut.GetTodoByTitle("Todo_new")
		assert.NoError(t, err)
		assert.Equal(t, want, got)

		want = Todo{Title: "Todo_updated", Completed: "false"}
		updatedTodo, err := sut.UpdateTodoTitle("Todo_new", "Todo_updated")
		assert.NoError(t, err)
		assert.Equal(t, want, updatedTodo)
		got, err = sut.GetTodoByTitle("Todo_updated")
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}
