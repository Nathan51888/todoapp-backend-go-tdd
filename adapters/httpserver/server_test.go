package httpserver_test

import (
	"encoding/json"
	"mytodoapp/adapters/httpserver"
	"mytodoapp/adapters/persistence/inmemory"
	"mytodoapp/domain/todo"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTodoServer(t *testing.T) {
	t.Run("can get todo by title", func(t *testing.T) {
		server := httpserver.NewTodoServer(&inmemory.InMemoryTodoStore{Todos: []todo.Todo{
			{Title: "Todo1", Completed: "false"},
		}})

		req := httptest.NewRequest(http.MethodGet, "/todo?title=Todo1", nil)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		var got todo.Todo
		json.NewDecoder(res.Body).Decode(&got)

		want := todo.Todo{Title: "Todo1", Completed: "false"}
		assert.Equal(t, want, got)
	})
}
