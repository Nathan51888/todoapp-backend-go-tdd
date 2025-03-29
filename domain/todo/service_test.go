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
	userId := uuid.New()
	userStore := &inmemory.InMemoryUserStore{Users: []user.User{
		{Id: userId},
	}}
	todoId := uuid.New()
	todoStore := &inmemory.InMemoryTodoStore{Todos: []todo.Todo{
		{Id: todoId, Title: "Todo1", Completed: false, UserId: userId},
	}}
	sut := todo.NewTodoService(todoStore, userStore)

	t.Run("can get todo by title", func(t *testing.T) {
		title := "Todo1"

		got, err := sut.GetTodoByTitle(userId, title)
		assert.NoError(t, err)
		want := todo.Todo{Id: todoId, Title: "Todo1", Completed: false, UserId: userId}
		assert.Equal(t, want, got)
	})
	t.Run("can get todo by id", func(t *testing.T) {
		id := todoId

		got, err := sut.GetTodoById(userId, id)
		assert.NoError(t, err)
		want := todo.Todo{Id: todoId, Title: "Todo1", Completed: false, UserId: userId}
		assert.Equal(t, want, got)
	})
}

func TestGetAllTodos(t *testing.T) {
	t.Run("can get all todos as slice", func(t *testing.T) {
		userId := uuid.New()
		userStore := &inmemory.InMemoryUserStore{Users: []user.User{
			{Id: userId},
		}}
		todoStore := &inmemory.InMemoryTodoStore{Todos: []todo.Todo{
			{Title: "Todo1", Completed: false, UserId: userId},
			{Title: "Todo2", Completed: true, UserId: userId},
			{Title: "Todo3", Completed: false, UserId: userId},
		}}
		sut := todo.NewTodoService(todoStore, userStore)
		want := []todo.Todo{
			{Title: "Todo1", Completed: false, UserId: userId},
			{Title: "Todo2", Completed: true, UserId: userId},
			{Title: "Todo3", Completed: false, UserId: userId},
		}

		got, err := sut.GetTodoAll(userId)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestCreateTodo(t *testing.T) {
	userId := uuid.New()
	userStore := &inmemory.InMemoryUserStore{Users: []user.User{
		{Id: userId},
	}}
	todoStore := &inmemory.InMemoryTodoStore{Todos: []todo.Todo{}}
	sut := todo.NewTodoService(todoStore, userStore)

	t.Run("can create todo with a todo object", func(t *testing.T) {
		todoToAdd := todo.Todo{Title: "Todo_new", Completed: true}

		got, err := sut.CreateTodo(userId, todoToAdd)
		assert.NoError(t, err)
		want := todo.Todo{Title: "Todo_new", Completed: true, UserId: userId}
		assert.Equal(t, want.Title, got.Title)
		assert.Equal(t, want.Completed, got.Completed)
		assert.Equal(t, want.UserId, got.UserId)
	})
	t.Run("can create todo with correct title", func(t *testing.T) {
		todoTitle := "Todo_new"
		todoToAdd := todo.Todo{Title: todoTitle}

		got, err := sut.CreateTodo(userId, todoToAdd)
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

func TestUpdateTodo(t *testing.T) {
	t.Run("can update todo's title by todo id", func(t *testing.T) {
		userId := uuid.New()
		userStore := &inmemory.InMemoryUserStore{Users: []user.User{
			{Id: userId},
		}}
		todoId := uuid.New()
		todoStore := &inmemory.InMemoryTodoStore{Todos: []todo.Todo{
			{Id: todoId, Title: "Todo_new", Completed: false, UserId: userId},
		}}
		sut := todo.NewTodoService(todoStore, userStore)
		want := todo.Todo{Id: todoId, Title: "Todo_updated", Completed: false, UserId: userId}

		got, err := sut.UpdateTodoTitle(userId, todoId, want.Title)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
	t.Run("can update todo's completed by todo id", func(t *testing.T) {
		userId := uuid.New()
		userStore := &inmemory.InMemoryUserStore{Users: []user.User{
			{Id: userId},
		}}
		todoId := uuid.New()
		todoStore := &inmemory.InMemoryTodoStore{Todos: []todo.Todo{
			{Id: todoId, Title: "Todo_new", Completed: false, UserId: userId},
		}}
		sut := todo.NewTodoService(todoStore, userStore)
		want := todo.Todo{Id: todoId, Title: "Todo_new", Completed: true, UserId: userId}

		got, err := sut.UpdateTodoStatus(userId, todoId, want.Completed)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
	t.Run("can update all fields on todo", func(t *testing.T) {
		userId := uuid.New()
		userStore := &inmemory.InMemoryUserStore{Users: []user.User{
			{Id: userId},
		}}
		todoId := uuid.New()
		todoStore := &inmemory.InMemoryTodoStore{Todos: []todo.Todo{
			{Id: todoId, Title: "Todo_new", Completed: false, UserId: userId},
		}}
		sut := todo.NewTodoService(todoStore, userStore)
		want := todo.Todo{Id: todoId, Title: "Todo_updated", Completed: true, UserId: userId}

		got, err := sut.UpdateTodoById(userId, todoId, want)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestDeleteTodo(t *testing.T) {
	t.Run("DELETE /todo: can delete todo by id", func(t *testing.T) {
		userId := uuid.New()
		userStore := &inmemory.InMemoryUserStore{Users: []user.User{
			{Id: userId},
		}}
		todoId := uuid.New()
		todoStore := &inmemory.InMemoryTodoStore{Todos: []todo.Todo{
			{Id: todoId, Title: "Delete_this", Completed: false, UserId: userId},
		}}
		sut := todo.NewTodoService(todoStore, userStore)

		got, err := sut.DeleteTodo(userId, todoId)
		assert.NoError(t, err)
		want := todo.Todo{Id: todoId, Title: "Delete_this", Completed: false, UserId: userId}
		assert.Equal(t, want, got)
	})
}
