package todo_test

import (
	"mytodoapp/adapters/persistence/inmemory"
	"mytodoapp/domain/todo"
	"mytodoapp/domain/user"
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
	sut := todo.NewTodoService(todoStore, userStore)

	t.Run("can get todo by title", func(t *testing.T) {
		title := "Todo1"

		got, err := sut.GetTodoByTitle(userId, title)
		assert.NoError(t, err)
		want := todo.Todo{Id: todoId, Title: "Todo1", Completed: false}
		assert.Equal(t, want, got)
	})
	t.Run("GET /todo: can get todo by id", func(t *testing.T) {
		id := todoId

		got, err := sut.GetTodoById(userId, id)
		assert.NoError(t, err)
		want := todo.Todo{Id: todoId, Title: "Todo1", Completed: false}
		assert.Equal(t, want, got)
	})
}

func TestGetAllTodos(t *testing.T) {
	t.Run("GET /todo: can get all todos as slice", func(t *testing.T) {
		todoStore := &inmemory.InMemoryTodoStore{Todos: []todo.Todo{
			{Title: "Todo1", Completed: false},
			{Title: "Todo2", Completed: true},
			{Title: "Todo3", Completed: false},
		}}
		userId := uuid.New()
		userStore := &inmemory.InMemoryUserStore{Users: []user.User{
			{Id: userId},
		}}
		sut := todo.NewTodoService(todoStore, userStore)
		want := []todo.Todo{
			{Title: "Todo1", Completed: false},
			{Title: "Todo2", Completed: true},
			{Title: "Todo3", Completed: false},
		}

		got, err := sut.GetTodoAll(userId)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestPOST(t *testing.T) {
	todoStore := &inmemory.InMemoryTodoStore{Todos: []todo.Todo{}}
	userId := uuid.New()
	userStore := &inmemory.InMemoryUserStore{Users: []user.User{
		{Id: userId},
	}}
	sut := todo.NewTodoService(todoStore, userStore)

	t.Run("can create todo with json", func(t *testing.T) {
		todoToAdd := todo.Todo{Title: "Todo_new", Completed: false}

		got, err := sut.CreateTodo(userId, todoToAdd)
		assert.NoError(t, err)
		want := todo.Todo{Title: "Todo_new", Completed: false}
		assert.Equal(t, want.Title, got.Title)
		assert.Equal(t, want.Completed, got.Completed)
	})
	t.Run("can create todo with query strings", func(t *testing.T) {
		todoTitle := "Todo_new"

		// TODO: take todo object as param
		got, err := sut.CreateTodo(userId, todo.Todo{Title: todoTitle})
		assert.NoError(t, err)
		want := todo.Todo{Title: "Todo_new", Completed: false}
		assert.Equal(t, want.Title, got.Title)
		assert.Equal(t, want.Completed, got.Completed)
	})
	t.Run("cannot create todo with empty title", func(t *testing.T) {
		todoTitle := ""

		_, err := sut.CreateTodo(userId, todo.Todo{Title: todoTitle})
		assert.Error(t, err)
	})
}

func TestPUT(t *testing.T) {
	t.Run("can update todo's title by todo id", func(t *testing.T) {
		todoId := uuid.New()
		todoStore := &inmemory.InMemoryTodoStore{Todos: []todo.Todo{
			{Id: todoId, Title: "Todo_new", Completed: false},
		}}
		userId := uuid.New()
		userStore := &inmemory.InMemoryUserStore{Users: []user.User{
			{Id: userId},
		}}
		sut := todo.NewTodoService(todoStore, userStore)
		want := todo.Todo{Id: todoId, Title: "Todo_updated", Completed: false}

		got, err := sut.UpdateTodoTitle(userId, todoId, want.Title)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
	t.Run("can update todo's completed by todo id", func(t *testing.T) {
		todoId := uuid.New()
		todoStore := &inmemory.InMemoryTodoStore{Todos: []todo.Todo{
			{Id: todoId, Title: "Todo_new", Completed: false},
		}}
		userId := uuid.New()
		userStore := &inmemory.InMemoryUserStore{Users: []user.User{
			{Id: userId},
		}}
		sut := todo.NewTodoService(todoStore, userStore)
		want := todo.Todo{Id: todoId, Title: "Todo_new", Completed: true}

		got, err := sut.UpdateTodoStatus(userId, todoId, want.Completed)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

// func TestDELETE(t *testing.T) {
// 	todoId := uuid.New()
// 	todoStore := &inmemory.InMemoryTodoStore{Todos: []todo.Todo{
// 		{Id: todoId, Title: "Todo1", Completed: false},
// 	}}
// 	userId := uuid.New()
// 	userStore := &inmemory.InMemoryUserStore{Users: []user.User{
// 		{Id: userId},
// 	}}
// 	handler := createTodoHandler(todoStore, userStore)
// 	token, err := auth.CreateAccessToken(userId.String())
// 	if err != nil {
// 		t.Fatalf("CreateJWT(): %v", err)
// 	}
//
// 	t.Run("DELETE /todo: cannot delete todo without auth header", func(t *testing.T) {
// 		req := httptest.NewRequest(http.MethodDelete, "/todo", nil)
// 		res := httptest.NewRecorder()
//
// 		handler.ServeHTTP(res, req)
//
// 		assert.Equal(t, http.StatusUnauthorized, res.Code)
// 	})
// 	t.Run("DELETE /todo: can delete todo by id", func(t *testing.T) {
// 		id := uuid.New()
// 		handler := createTodoHandler(&inmemory.InMemoryTodoStore{Todos: []todo.Todo{
// 			{Id: id, Title: "Delete_this", Completed: false},
// 		}}, userStore)
//
// 		req := httptest.NewRequest(http.MethodDelete, "/todo?id="+id.String(), nil)
// 		req.Header.Add("Authorization", token)
// 		res := httptest.NewRecorder()
//
// 		handler.ServeHTTP(res, req)
//
// 		want := todo.Todo{Id: id, Title: "Delete_this", Completed: false}
// 		var got todo.Todo
// 		json.NewDecoder(res.Body).Decode(&got)
// 		assert.Equal(t, want, got)
// 	})
// }
