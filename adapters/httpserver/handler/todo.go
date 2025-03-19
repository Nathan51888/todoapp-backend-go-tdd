package handler

import (
	"encoding/json"
	"log"
	"mytodoapp/adapters/auth"
	"mytodoapp/domain/todo"
	"mytodoapp/domain/user"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type TodoHandler struct {
	todoStore todo.TodoStore
	userStore user.UserStore
}

func NewTodoHandler(mux *http.ServeMux, todoStore todo.TodoStore, userStore user.UserStore) {
	handler := &TodoHandler{
		todoStore: todoStore,
		userStore: userStore,
	}
	mux.HandleFunc("GET /todo", auth.WithJWTAuth(handler.GetTodo, userStore))
	mux.HandleFunc("POST /todo", auth.WithJWTAuth(handler.CreateTodo, userStore))
	mux.HandleFunc("PUT /todo", auth.WithJWTAuth(handler.UpdateTodo, userStore))
	mux.HandleFunc("DELETE /todo", auth.WithJWTAuth(handler.DeleteTodo, userStore))
}

func (t *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	id := r.URL.Query().Get("id")

	if title != "" {
		t.GetTodoByTitle(w, r, title)
		return
	}

	if id != "" {
		t.GetTodoById(w, r, id)
		return
	}

	t.GetTodoAll(w, r)
}

func (t *TodoHandler) GetTodoByTitle(w http.ResponseWriter, r *http.Request, title string) {
	result, err := t.todoStore.GetTodoByTitle(title)
	if err != nil {
		log.Printf("Error GetTodoByTitle(): %v", err)
	}

	json.NewEncoder(w).Encode(&result)
}

func (t *TodoHandler) GetTodoById(w http.ResponseWriter, r *http.Request, idString string) {
	id, err := uuid.Parse(idString)
	if err != nil {
		log.Printf("Error converting string to int: %v", err)
	}

	result, err := t.todoStore.GetTodoById(id)
	if err != nil {
		log.Printf("Error GetTodoById(): %v", err)
	}

	json.NewEncoder(w).Encode(result)
}

func (t *TodoHandler) GetTodoAll(w http.ResponseWriter, r *http.Request) {
	result, err := t.todoStore.GetTodoAll()
	if err != nil {
		log.Printf("Error GetTodoAll(): %v", err)
	}

	json.NewEncoder(w).Encode(&result)
}

func (t *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")

	var body todo.Todo
	json.NewDecoder(r.Body).Decode(&body)
	if body.Title != "" {
		result, err := t.todoStore.CreateTodo(body.Title)
		if err != nil {
			log.Printf("Error CreateTodo(): %v", err)
		}

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&result)
	}

	if title == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if title != "" {
		result, err := t.todoStore.CreateTodo(title)
		if err != nil {
			log.Printf("Error CreateTodo(): %v", err)
		}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&result)
		return
	}
}

func (t *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	title := r.URL.Query().Get("title")
	completed := r.URL.Query().Get("completed")

	if title != "" {
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, err := uuid.Parse(id)
		if err != nil {
			log.Printf("Error parsing id string: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		result, err := t.todoStore.UpdateTodoTitle(id, title)
		if err != nil {
			log.Printf("Error UpdateTodoById(): %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(&result)
		return
	}

	if completed != "" {
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, err := uuid.Parse(id)
		if err != nil {
			log.Printf("Error parsing id string: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		completed, err := strconv.ParseBool(completed)
		if err != nil {
			log.Printf("Error parsing completed string to bool: %v", err)
		}
		result, err := t.todoStore.UpdateTodoStatus(id, completed)
		if err != nil {
			log.Printf("Error UpdateTodoById(): %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(&result)
		return
	}

	// json body
	var body todo.Todo
	json.NewDecoder(r.Body).Decode(&body)
	result, err := t.todoStore.UpdateTodoById(body.Id, body)
	if err != nil {
		log.Printf("Error UpdateTodoById(): %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(&result)
}

func (t *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("DeleteTodo: %v", err)
		return
	}
	result, err := t.todoStore.DeleteTodoById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("DeleteTodo: %v", err)
	}

	json.NewEncoder(w).Encode(result)
}
