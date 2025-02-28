package specifications

import (
	"mytodoapp/domain/todo"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TodoDriver interface {
	GetTodoByTitle(title string) (todo.Todo, error)
}

func TodoSpecification(t testing.TB, driver TodoDriver) {
	got, err := driver.GetTodoByTitle("Todo1")
	want := todo.Todo{Title: "Todo1", Completed: "false"}
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
