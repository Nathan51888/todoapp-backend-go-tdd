package specifications

import (
	"mytodoapp/domain/todo"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TodoDriver interface {
	GetTodoByTitle(title string) (todo.Todo, error)
	CreateTodo(title string) (todo.Todo, error)
}

func TodoSpecification(t testing.TB, driver TodoDriver) {
	want := todo.Todo{Title: "Todo_new", Completed: "false"}
	got, err := driver.CreateTodo("Todo_new")
	assert.NoError(t, err)
	assert.Equal(t, want, got)

	got, err = driver.GetTodoByTitle("Todo_new")
	want = todo.Todo{Title: "Todo_new", Completed: "false"}
	assert.NoError(t, err)
	assert.Equal(t, want, got)

	got, err = driver.GetTodoByTitle("Todo_new")
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
