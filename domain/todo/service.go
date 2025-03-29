package todo

import (
	"errors"
	"fmt"
	"mytodoapp/domain/user"

	"github.com/google/uuid"
)

type TodoService struct {
	todoStore TodoStore
	userStore user.UserStore
}

func NewTodoService(todoStore TodoStore, userStore user.UserStore) *TodoService {
	return &TodoService{todoStore, userStore}
}

func (t *TodoService) GetTodoByTitle(userId uuid.UUID, title string) (Todo, error) {
	result, err := t.todoStore.GetTodoByTitle(userId, title)
	if errors.Is(err, ErrTodoNotFound) {
		return Todo{}, NewError(ErrNotFound, err)
	}
	if err != nil {
		return Todo{}, NewError(ErrInternalError, fmt.Errorf("GetTodoByTitle(): %v", err))
	}

	return result, nil
}

func (t *TodoService) GetTodoById(userId uuid.UUID, todoId uuid.UUID) (Todo, error) {
	result, err := t.todoStore.GetTodoById(userId, todoId)
	if errors.Is(err, ErrTodoNotFound) {
		return Todo{}, NewError(ErrNotFound, err)
	}
	if err != nil {
		return Todo{}, NewError(ErrInternalError, fmt.Errorf("GetTodoById(): %v", err))
	}

	return result, nil
}

func (t *TodoService) GetTodoAll(userId uuid.UUID) ([]Todo, error) {
	result, err := t.todoStore.GetTodoAll(userId)
	if errors.Is(err, ErrTodoNotFound) {
		return nil, NewError(ErrNotFound, err)
	}
	if err != nil {
		return nil, NewError(ErrInternalError, fmt.Errorf("GetTodoAll(): %v", err))
	}

	return result, nil
}

func (t *TodoService) CreateTodo(userId uuid.UUID, todoToAdd Todo) (Todo, error) {
	if userId == uuid.Nil {
		return Todo{}, NewError(ErrBadRequest, ErrUserIdEmpty)
	}
	if todoToAdd.Title == "" {
		return Todo{}, NewError(ErrBadRequest, ErrTodoTitleEmpty)
	}

	result, err := t.todoStore.CreateTodo(userId, todoToAdd)
	if errors.Is(err, ErrTodoTitleEmpty) {
		return Todo{}, NewError(ErrBadRequest, fmt.Errorf("CreateTodo: %w", err))
	}
	if err != nil {
		return Todo{}, NewError(ErrInternalError, fmt.Errorf("CreateTodo: %w", err))
	}

	return result, nil
}

func (t *TodoService) UpdateTodoTitle(userId uuid.UUID, todoId uuid.UUID, title string) (Todo, error) {
	if userId == uuid.Nil {
		return Todo{}, NewError(ErrBadRequest, ErrUserIdEmpty)
	}
	if todoId == uuid.Nil {
		return Todo{}, NewError(ErrBadRequest, ErrTodoIdEmpty)
	}

	result, err := t.todoStore.UpdateTodoTitle(userId, todoId, title)
	if errors.Is(err, ErrTodoNotFound) {
		return Todo{}, NewError(ErrNotFound, err)
	}
	if err != nil {
		return Todo{}, NewError(ErrInternalError, err)
	}

	return result, nil
}

func (t *TodoService) UpdateTodoById(userId uuid.UUID, todoId uuid.UUID, changedTodo Todo) (Todo, error) {
	if userId == uuid.Nil {
		return Todo{}, NewError(ErrBadRequest, ErrUserIdEmpty)
	}
	if todoId == uuid.Nil {
		return Todo{}, NewError(ErrBadRequest, ErrTodoIdEmpty)
	}

	// Full update
	result, err := t.todoStore.UpdateTodoById(userId, todoId, changedTodo)
	if errors.Is(err, ErrTodoNotFound) {
		return Todo{}, NewError(ErrNotFound, err)
	}
	if err != nil {
		return Todo{}, NewError(ErrInternalError, err)
	}

	return result, nil
}

func (t *TodoService) UpdateTodoStatus(userId uuid.UUID, todoId uuid.UUID, completed bool) (Todo, error) {
	if userId == uuid.Nil {
		return Todo{}, NewError(ErrBadRequest, ErrUserIdEmpty)
	}
	if todoId == uuid.Nil {
		return Todo{}, NewError(ErrBadRequest, ErrTodoTitleEmpty)
	}

	result, err := t.todoStore.UpdateTodoStatus(userId, todoId, completed)
	if errors.Is(err, ErrTodoNotFound) {
		return Todo{}, NewError(ErrNotFound, err)
	}
	if err != nil {
		return Todo{}, NewError(ErrInternalError, fmt.Errorf("UpdateTodoStatus(): %v", err))
	}

	return result, nil
}

func (t *TodoService) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("DeleteTodo: %v", err)
		return
	}

	userId, err := auth.GetUserIdFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("auth.GetUserIdFromContext: %v", err)
		return
	}

	result, err := t.todoStore.DeleteTodoById(userId, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("DeleteTodo: %v", err)
		return
	}

	json.NewEncoder(w).Encode(result)
}
