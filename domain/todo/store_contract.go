package todo

import (
	"mytodoapp/domain/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

// make user store a dependency
type TodoStoreContract struct {
	NewTodoStore func() (TodoStore, error)
	UserStore    user.UserStore
}

func (c TodoStoreContract) Test(t *testing.T) {
	_, err := c.UserStore.CreateUser("test@email.com", "test")
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	got, err := c.UserStore.GetUserByEmail("test@email.com")
	if err != nil {
		t.Fatalf("failed to get user by emai: %v", err)
	}
	userId := got.Id

	sut, err := c.NewTodoStore()
	if err != nil {
		t.Fatalf("Error creating todo store: %v\n", err)
	}

	t.Run("GetTodoAll(): can get all todos as slice from database", func(t *testing.T) {
		want := []Todo{
			{Title: "Todo1", Completed: false, UserId: userId},
			{Title: "Todo2", Completed: false, UserId: userId},
			{Title: "Todo3", Completed: false, UserId: userId},
		}
		sut.CreateTodoWithTitle(userId, "Todo1")
		sut.CreateTodoWithTitle(userId, "Todo2")
		sut.CreateTodoWithTitle(userId, "Todo3")

		got, err := sut.GetTodoAll(userId)
		assert.NoError(t, err)
		for index, item := range got {
			if item.Title != want[index].Title {
				t.Error("title not equal")
			}
			if item.Completed != want[index].Completed {
				t.Error("completed not equal")
			}
			if item.UserId != want[index].UserId {
				t.Error("userId not equal")
			}
		}
	})
	t.Run("can create, get, update todo's title and status by title from database", func(t *testing.T) {
		want := Todo{Title: "Todo", Completed: false, UserId: userId}
		newTodo, err := sut.CreateTodo(userId, want)
		assert.NoError(t, err)
		assert.Equal(t, want.Title, newTodo.Title)
		assert.Equal(t, want.Completed, newTodo.Completed)
		assert.Equal(t, want.UserId, newTodo.UserId)

		want = Todo{Id: newTodo.Id, Title: "Todo", Completed: false, UserId: userId}
		got, err := sut.GetTodoById(userId, newTodo.Id)
		assert.NoError(t, err)
		assert.Equal(t, want, got, "GetTodoById()")

		want = Todo{Id: got.Id, Title: "Todo_updated", Completed: got.Completed, UserId: userId}
		updatedTodo, err := sut.UpdateTodoTitle(userId, want.Id, "Todo_updated")
		assert.NoError(t, err)
		assert.Equal(t, want, updatedTodo, "UpdateTodoTitle()")
		got, err = sut.GetTodoByTitle(userId, "Todo_updated")
		assert.NoError(t, err)
		assert.Equal(t, want, got, "GetTodoByTitle()")

		want = Todo{Id: got.Id, Title: "Todo_updated", Completed: true, UserId: userId}
		updatedTodo, err = sut.UpdateTodoStatus(userId, want.Id, true)
		assert.NoError(t, err)
		assert.Equal(t, want, updatedTodo, "UpdateTodoStatus()")
		got, err = sut.GetTodoByTitle(userId, "Todo_updated")
		assert.NoError(t, err)
		assert.Equal(t, want, got, "GetTodoByTitle()")
	})
	t.Run("CreateTodo(): can create todo with correct userId when its not provided in todo object", func(t *testing.T) {
		want := Todo{Title: "Todo", Completed: false, UserId: userId}
		newTodo, err := sut.CreateTodo(userId, want)
		assert.NoError(t, err)
		assert.Equal(t, want.Title, newTodo.Title)
		assert.Equal(t, want.Completed, newTodo.Completed)
		assert.Equal(t, want.UserId, newTodo.UserId)

		want = Todo{Id: newTodo.Id, Title: "Todo", Completed: false, UserId: userId}
		got, err := sut.GetTodoById(userId, newTodo.Id)
		assert.NoError(t, err)
		assert.Equal(t, want, got, "GetTodoById()")
	})
	t.Run("CreateTodo(): cannot create todo with empty title", func(t *testing.T) {
		want := Todo{Title: "", Completed: false, UserId: userId}
		newTodo, err := sut.CreateTodo(userId, want)
		assert.Error(t, err)
		assert.Equal(t, Todo{}, newTodo)

		want = Todo{Id: newTodo.Id, Title: want.Title, Completed: false, UserId: userId}
		got, err := sut.GetTodoById(userId, newTodo.Id)
		assert.Error(t, err)
		assert.NotEqual(t, want, got, "GetTodoById()")
	})
	t.Run("CreateTodoWithTitle(): cannot create todo with empty title", func(t *testing.T) {
		want := Todo{Title: "", Completed: false, UserId: userId}
		newTodo, err := sut.CreateTodoWithTitle(userId, want.Title)
		assert.Error(t, err)
		assert.Equal(t, Todo{}, newTodo)

		want = Todo{Id: newTodo.Id, Title: want.Title, Completed: false, UserId: userId}
		got, err := sut.GetTodoById(userId, newTodo.Id)
		assert.Error(t, err)
		assert.NotEqual(t, want, got, "GetTodoById()")
	})
	t.Run("can update todo by id", func(t *testing.T) {
		want := Todo{Title: "Todo_new", Completed: false}
		newTodo, err := sut.CreateTodo(userId, want)
		assert.NoError(t, err)
		assert.Equal(t, want.Title, newTodo.Title)
		assert.Equal(t, want.Completed, newTodo.Completed)
		got, err := sut.GetTodoById(userId, newTodo.Id)
		assert.NoError(t, err)
		assert.Equal(t, newTodo, got, "GetTodoById()")

		want = Todo{Id: got.Id, Title: "Todo_updated", Completed: true, UserId: userId}
		updatedTodo, err := sut.UpdateTodoById(userId, want.Id, want)
		assert.NoError(t, err)
		assert.Equal(t, want, updatedTodo, "UpdateTodoById()")
		got, err = sut.GetTodoById(userId, updatedTodo.Id)
		assert.NoError(t, err)
		assert.Equal(t, want, got, "UpdateTodoById()")
	})
	t.Run("can delete todo by id", func(t *testing.T) {
		// TODO: dry it with function
		want := Todo{Title: "Delete_this", Completed: false}
		newTodo, err := sut.CreateTodoWithTitle(userId, "Delete_this")
		assert.NoError(t, err)
		assert.Equal(t, want.Title, newTodo.Title)
		assert.Equal(t, want.Completed, newTodo.Completed)
		got, err := sut.GetTodoById(userId, newTodo.Id)
		assert.NoError(t, err)
		assert.Equal(t, newTodo, got, "GetTodoById()")

		want = Todo{Id: got.Id, Title: "Delete_this", Completed: false, UserId: userId}
		deletedTodo, err := sut.DeleteTodoById(userId, want.Id)
		assert.NoError(t, err, "DeleteTodoById()")
		assert.Equal(t, want, deletedTodo, "DeleteTodoById()")
		got, err = sut.GetTodoById(userId, want.Id)
		assert.Error(t, err, "GetTodoById()")
		assert.NotEqual(t, want, got, "GetTodoById()")
	})
}
