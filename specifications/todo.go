package specifications

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TodoDriver interface {
	GetTodoByTitle(title string) (string, error)
}

func TodoSpecification(t testing.TB, driver TodoDriver) {
	got, err := driver.GetTodoByTitle("Todo1")
	assert.NoError(t, err)
	assert.Equal(t, "Todo1", got)
}
