package todo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TodoStoreContract struct {
	NewTodoStore func() (TodoStore, error)
}

func (c TodoStoreContract) Test(t *testing.T) {
	t.Run("can create, get, update todo's title and status by title from database", func(t *testing.T) {
		sut, err := c.NewTodoStore()
		if err != nil {
			t.Fatalf("Error creating todo store: %v\n", err)
		}

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

		want = Todo{Title: "Todo_updated", Completed: "true"}
		updatedTodo, err = sut.UpdateTodoStatus("Todo_updated", "true")
		assert.NoError(t, err, "UpdateTodoStatus()")
		assert.Equal(t, want, updatedTodo, "UpdateTodoStatus()")
		got, err = sut.GetTodoByTitle("Todo_updated")
		assert.NoError(t, err, "GetTodoByTitle()")
		assert.Equal(t, want, got, "GetTodoByTitle()")
	})
	t.Run("can get all todos as slice from database", func(t *testing.T) {
		sut, err := c.NewTodoStore()
		if err != nil {
			t.Fatalf("Error creating todo store: %v\n", err)
		}

		want := []Todo{
			{Title: "Todo1", Completed: "false"},
			{Title: "Todo2", Completed: "false"},
			{Title: "Todo3", Completed: "false"},
		}
		sut.CreateTodo("Todo1")
		sut.CreateTodo("Todo2")
		sut.CreateTodo("Todo3")

		got, err := sut.GetTodoAll()
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}
