package httpserver_test

import (
	"bytes"
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
			{Title: "Todo1", Completed: false},
		}})

		req := httptest.NewRequest(http.MethodGet, "/todo?title=Todo1", nil)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		var got todo.Todo
		json.NewDecoder(res.Body).Decode(&got)

		want := todo.Todo{Title: "Todo1", Completed: false}
		assert.Equal(t, want, got)
	})
	t.Run("GET /todo: can get todo by id", func(t *testing.T) {
		server := httpserver.NewTodoServer(&inmemory.InMemoryTodoStore{Todos: []todo.Todo{
			{Id: 3, Title: "Todo1", Completed: false},
		}})

		req := httptest.NewRequest(http.MethodGet, "/todo?id=3", nil)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		var got todo.Todo
		json.NewDecoder(res.Body).Decode(&got)

		want := todo.Todo{Id: 3, Title: "Todo1", Completed: false}
		assert.Equal(t, want, got)
	})
	t.Run("POST /todo: can create and get todo by title", func(t *testing.T) {
		server := httpserver.NewTodoServer(&inmemory.InMemoryTodoStore{})

		body := todo.Todo{Title: "Todo_new", Completed: false}
		payloadBuf := new(bytes.Buffer)
		json.NewEncoder(payloadBuf).Encode(body)
		req := httptest.NewRequest(http.MethodPost, "/todo", payloadBuf)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		var got todo.Todo
		json.NewDecoder(res.Body).Decode(&got)
		want := todo.Todo{Title: "Todo_new", Completed: false}
		assert.Equal(t, want.Title, got.Title)
		assert.Equal(t, want.Completed, got.Completed)
	})
	t.Run("POST /todo?title: can create todo with query options", func(t *testing.T) {
		server := httpserver.NewTodoServer(&inmemory.InMemoryTodoStore{})

		req := httptest.NewRequest(http.MethodPost, "/todo?title=Todo_new", nil)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		var got todo.Todo
		json.NewDecoder(res.Body).Decode(&got)
		want := todo.Todo{Title: "Todo_new", Completed: false}
		assert.Equal(t, want.Title, got.Title)
		assert.Equal(t, want.Completed, got.Completed)
	})
	t.Run("PUT /todo: can update todo title by todo id", func(t *testing.T) {
		server := httpserver.NewTodoServer(&inmemory.InMemoryTodoStore{Todos: []todo.Todo{
			{Id: 5, Title: "Todo_new", Completed: false},
		}})

		want := todo.Todo{Id: 5, Title: "Todo_updated", Completed: false}

		payloadBuf := new(bytes.Buffer)
		json.NewEncoder(payloadBuf).Encode(want)
		req := httptest.NewRequest(http.MethodPut, "/todo", payloadBuf)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		var got todo.Todo
		json.NewDecoder(res.Body).Decode(&got)
		assert.Equal(t, want, got)
	})
	t.Run("PUT /todo: can update todo's status by id with body", func(t *testing.T) {
		server := httpserver.NewTodoServer(&inmemory.InMemoryTodoStore{Todos: []todo.Todo{
			{Id: 5, Title: "Todo_new", Completed: false},
		}})

		want := todo.Todo{Id: 5, Title: "Todo_new", Completed: true}

		payloadBuf := new(bytes.Buffer)
		json.NewEncoder(payloadBuf).Encode(want)
		req := httptest.NewRequest(http.MethodPut, "/todo", payloadBuf)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		var got todo.Todo
		json.NewDecoder(res.Body).Decode(&got)
		assert.Equal(t, want, got)
	})
	t.Run("GET /todo: can get all todos as slice", func(t *testing.T) {
		want := []todo.Todo{
			{Title: "Todo1", Completed: false},
			{Title: "Todo2", Completed: true},
			{Title: "Todo3", Completed: false},
		}
		server := httpserver.NewTodoServer(&inmemory.InMemoryTodoStore{Todos: want})

		req := httptest.NewRequest(http.MethodGet, "/todo", nil)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		var got []todo.Todo
		json.NewDecoder(res.Body).Decode(&got)
		assert.Equal(t, want, got)
	})
	t.Run("DELETE /todo: can delete todo by id", func(t *testing.T) {
		server := httpserver.NewTodoServer(&inmemory.InMemoryTodoStore{Todos: []todo.Todo{
			{Id: 8, Title: "Delete_this", Completed: false},
		}})

		req := httptest.NewRequest(http.MethodDelete, "/todo?id=8", nil)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		want := todo.Todo{Id: 8, Title: "Delete_this", Completed: false}
		var got todo.Todo
		json.NewDecoder(res.Body).Decode(&got)
		assert.Equal(t, want, got)
	})
}
