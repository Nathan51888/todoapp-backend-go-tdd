package handler_test

import (
	"bytes"
	"encoding/json"
	"mytodoapp/adapters/auth"
	"mytodoapp/adapters/httpserver/handler"
	"mytodoapp/adapters/persistence/inmemory"
	"mytodoapp/domain/todo"
	"mytodoapp/domain/user"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetTodo(t *testing.T) {
	todoId := uuid.New()
	todoStore := &inmemory.InMemoryTodoStore{Todos: []todo.Todo{
		{Id: todoId, Title: "Todo1", Completed: false},
	}}
	userId := uuid.New()
	userStore := &inmemory.InMemoryUserStore{Users: []user.User{
		{Id: userId},
	}}
	handler := createTodoHandler(todoStore, userStore)
	token, err := auth.CreateJWT([]byte(auth.JWTSecret), userId.String())
	if err != nil {
		t.Fatalf("CreateJWT(): %v", err)
	}

	t.Run("GET /todo: cannot get todo without auth header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/todo", nil)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})
	t.Run("GET /todo: can get todo with auth header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/todo", nil)
		req.Header.Add("Authorization", token)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})
	t.Run("GET /todo: can get todo by title", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/todo?title=Todo1", nil)
		req.Header.Add("Authorization", token)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		var got todo.Todo
		json.NewDecoder(res.Body).Decode(&got)

		want := todo.Todo{Id: todoId, Title: "Todo1", Completed: false}
		assert.Equal(t, want, got)
	})
	t.Run("GET /todo: can get todo by id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/todo?id="+todoId.String(), nil)
		req.Header.Add("Authorization", token)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		var got todo.Todo
		json.NewDecoder(res.Body).Decode(&got)

		want := todo.Todo{Id: todoId, Title: "Todo1", Completed: false}
		assert.Equal(t, want, got)
	})
}

func TestGetAllTodos(t *testing.T) {
	todoStore := &inmemory.InMemoryTodoStore{Todos: []todo.Todo{
		{Title: "Todo1", Completed: false},
		{Title: "Todo2", Completed: true},
		{Title: "Todo3", Completed: false},
	}}
	userId := uuid.New()
	userStore := &inmemory.InMemoryUserStore{Users: []user.User{
		{Id: userId},
	}}
	handler := createTodoHandler(todoStore, userStore)
	token, err := auth.CreateJWT([]byte(auth.JWTSecret), userId.String())
	if err != nil {
		t.Fatalf("CreateJWT(): %v", err)
	}

	t.Run("GET /todo: cannot get todo without auth header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/todo", nil)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})
	t.Run("GET /todo: can get todo with auth header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/todo", nil)
		req.Header.Add("Authorization", token)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})
	t.Run("GET /todo: can get all todos as slice", func(t *testing.T) {
		want := []todo.Todo{
			{Title: "Todo1", Completed: false},
			{Title: "Todo2", Completed: true},
			{Title: "Todo3", Completed: false},
		}

		req := httptest.NewRequest(http.MethodGet, "/todo", nil)
		req.Header.Add("Authorization", token)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		var got []todo.Todo
		json.NewDecoder(res.Body).Decode(&got)
		assert.Equal(t, want, got)
	})
}

func TestPOST(t *testing.T) {
	todoStore := &inmemory.InMemoryTodoStore{Todos: []todo.Todo{}}
	userId := uuid.New()
	userStore := &inmemory.InMemoryUserStore{Users: []user.User{
		{Id: userId},
	}}
	handler := createTodoHandler(todoStore, userStore)
	token, err := auth.CreateJWT([]byte(auth.JWTSecret), userId.String())
	if err != nil {
		t.Fatalf("CreateJWT(): %v", err)
	}

	t.Run("POST /todo: cannot create todo without auth header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/todo", nil)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})
	t.Run("POST /todo: can create todo with json", func(t *testing.T) {
		body := todo.Todo{Title: "Todo_new", Completed: false}
		payloadBuf := new(bytes.Buffer)
		json.NewEncoder(payloadBuf).Encode(body)
		req := httptest.NewRequest(http.MethodPost, "/todo", payloadBuf)
		req.Header.Add("Authorization", token)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		var got todo.Todo
		json.NewDecoder(res.Body).Decode(&got)
		want := todo.Todo{Title: "Todo_new", Completed: false}
		assert.Equal(t, want.Title, got.Title)
		assert.Equal(t, want.Completed, got.Completed)
		assert.Equal(t, http.StatusCreated, res.Code)
	})
	t.Run("POST /todo?title: can create todo with query strings", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/todo?title=Todo_new", nil)
		req.Header.Add("Authorization", token)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		var got todo.Todo
		json.NewDecoder(res.Body).Decode(&got)
		want := todo.Todo{Title: "Todo_new", Completed: false}
		assert.Equal(t, want.Title, got.Title)
		assert.Equal(t, want.Completed, got.Completed)
		assert.Equal(t, http.StatusCreated, res.Code)
	})
	t.Run("POST /todo?title: cannot create todo with empty title", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/todo?title=", nil)
		req.Header.Add("Authorization", token)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
	})
}

func TestPUT(t *testing.T) {
	todoId := uuid.New()
	todoStore := &inmemory.InMemoryTodoStore{Todos: []todo.Todo{
		{Id: todoId, Title: "Todo_new", Completed: false},
	}}
	userId := uuid.New()
	userStore := &inmemory.InMemoryUserStore{Users: []user.User{
		{Id: userId},
	}}
	handler := createTodoHandler(todoStore, userStore)
	token, err := auth.CreateJWT([]byte(auth.JWTSecret), userId.String())
	if err != nil {
		t.Fatalf("CreateJWT(): %v", err)
	}

	t.Run("PUT /todo: cannot update todo without auth header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/todo", nil)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})
	t.Run("PUT /todo: can update todo with auth header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/todo", nil)
		req.Header.Add("Authorization", token)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})
	t.Run("PUT /todo: can update todo title by todo id", func(t *testing.T) {
		id := uuid.New()
		handler := createTodoHandler(&inmemory.InMemoryTodoStore{Todos: []todo.Todo{
			{Id: id, Title: "Todo_new", Completed: false},
		}}, userStore)

		want := todo.Todo{Id: id, Title: "Todo_updated", Completed: false}

		payloadBuf := new(bytes.Buffer)
		json.NewEncoder(payloadBuf).Encode(want)
		req := httptest.NewRequest(http.MethodPut, "/todo", payloadBuf)
		req.Header.Add("Authorization", token)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		var got todo.Todo
		json.NewDecoder(res.Body).Decode(&got)
		assert.Equal(t, want, got)
	})
	t.Run("PUT /todo: can update todo's status by id with body", func(t *testing.T) {
		id := uuid.New()
		handler := createTodoHandler(&inmemory.InMemoryTodoStore{Todos: []todo.Todo{
			{Id: id, Title: "Todo_new", Completed: false},
		}}, userStore)

		want := todo.Todo{Id: id, Title: "Todo_new", Completed: true}

		payloadBuf := new(bytes.Buffer)
		json.NewEncoder(payloadBuf).Encode(want)
		req := httptest.NewRequest(http.MethodPut, "/todo", payloadBuf)
		req.Header.Add("Authorization", token)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		var got todo.Todo
		json.NewDecoder(res.Body).Decode(&got)
		assert.Equal(t, want, got)
	})
}

func TestDELETE(t *testing.T) {
	todoId := uuid.New()
	todoStore := &inmemory.InMemoryTodoStore{Todos: []todo.Todo{
		{Id: todoId, Title: "Todo1", Completed: false},
	}}
	userId := uuid.New()
	userStore := &inmemory.InMemoryUserStore{Users: []user.User{
		{Id: userId},
	}}
	handler := createTodoHandler(todoStore, userStore)
	token, err := auth.CreateJWT([]byte(auth.JWTSecret), userId.String())
	if err != nil {
		t.Fatalf("CreateJWT(): %v", err)
	}

	t.Run("DELETE /todo: cannot delete todo without auth header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/todo", nil)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})
	t.Run("DELETE /todo: can delete todo by id", func(t *testing.T) {
		id := uuid.New()
		handler := createTodoHandler(&inmemory.InMemoryTodoStore{Todos: []todo.Todo{
			{Id: id, Title: "Delete_this", Completed: false},
		}}, userStore)

		req := httptest.NewRequest(http.MethodDelete, "/todo?id="+id.String(), nil)
		req.Header.Add("Authorization", token)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		want := todo.Todo{Id: id, Title: "Delete_this", Completed: false}
		var got todo.Todo
		json.NewDecoder(res.Body).Decode(&got)
		assert.Equal(t, want, got)
	})
}

func createTodoHandler(todoStore todo.TodoStore, userStore user.UserStore) *http.ServeMux {
	mux := http.NewServeMux()
	handler.NewTodoHandler(mux, todoStore, userStore)
	return mux
}
