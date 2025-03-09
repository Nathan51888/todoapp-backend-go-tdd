package specifications

import (
	"mytodoapp/domain/todo"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TodoDriver interface {
	GetTodoById(id int) (todo.Todo, error)
	CreateTodo(title string) (todo.Todo, error)
}

func TodoSpecification(t testing.TB, driver TodoDriver) {
	want := todo.Todo{Title: "Todo_new", Completed: false}
	got, err := driver.CreateTodo("Todo_new")
	assert.NoError(t, err)
	assert.Equal(t, want.Title, got.Title)
	assert.Equal(t, want.Completed, got.Completed)

	want = todo.Todo{Id: got.Id, Title: "Todo_new", Completed: false}
	got, err = driver.GetTodoById(want.Id)
	assert.NoError(t, err)
	assert.Equal(t, want, got, "driver.GetTodoById()")
}
