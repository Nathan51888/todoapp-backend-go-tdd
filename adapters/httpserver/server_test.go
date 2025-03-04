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
	t.Run("GET /todo: can get todo by title", func(t *testing.T) {
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
	t.Run("POST /todo: can create and get todo by title", func(t *testing.T) {
		server := httpserver.NewTodoServer(&inmemory.InMemoryTodoStore{})

		req := httptest.NewRequest(http.MethodPost, "/todo?title=Todo_new", nil)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		var got todo.Todo
		json.NewDecoder(res.Body).Decode(&got)

		want := todo.Todo{Title: "Todo_new", Completed: "false"}
		assert.Equal(t, want, got)
	})
	t.Run("PUT /todo: can update todo by title", func(t *testing.T) {
		server := httpserver.NewTodoServer(&inmemory.InMemoryTodoStore{Todos: []todo.Todo{
			{Title: "Todo_new", Completed: "false"},
		}})

		req := httptest.NewRequest(http.MethodPut, "/todo?target=Todo_new&title=Todo_updated", nil)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		var got todo.Todo
		json.NewDecoder(res.Body).Decode(&got)

		want := todo.Todo{Title: "Todo_updated", Completed: "false"}
		assert.Equal(t, want, got)
	})
}
