package handler

import (
	"encoding/json"
	"log"
	"mytodoapp/adapters/auth"
	"mytodoapp/adapters/httpserver/middleware"
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
	mux.HandleFunc("GET /todo", middleware.WithJWTAuth(handler.GetTodo, userStore))
	mux.HandleFunc("POST /todo", middleware.WithJWTAuth(handler.CreateTodo, userStore))
	mux.HandleFunc("PUT /todo", middleware.WithJWTAuth(handler.UpdateTodo, userStore))
	mux.HandleFunc("DELETE /todo", middleware.WithJWTAuth(handler.DeleteTodo, userStore))
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
	userId, err := auth.GetUserIdFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("auth.GetUserIdFromContext(): %v", err)
		return
	}
	result, err := t.todoStore.GetTodoByTitle(userId, title)
	if err != nil {
		log.Printf("Error GetTodoByTitle(): %v", err)
	}

	json.NewEncoder(w).Encode(&result)
}

func (t *TodoHandler) GetTodoById(w http.ResponseWriter, r *http.Request, idString string) {
	id, err := uuid.Parse(idString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error converting string to int: %v", err)
		return
	}

	userId, err := auth.GetUserIdFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("auth.GetUserIdFromContext(): %v", err)
		return
	}

	result, err := t.todoStore.GetTodoById(userId, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error GetTodoById(): %v", err)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (t *TodoHandler) GetTodoAll(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserIdFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("auth.GetUserIdFromContext(): %v", err)
		return
	}

	result, err := t.todoStore.GetTodoAll(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error GetTodoAll(): %v", err)
		return
	}

	json.NewEncoder(w).Encode(&result)
}

func (t *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	userId, err := auth.GetUserIdFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("auth.GetUserIdFromContext: %v", err)
		return
	}

	if title != "" {
		result, err := t.todoStore.CreateTodo(userId, title)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error CreateTodo(): %v", err)
			return
		}

		log.Printf("Todo created from query: %v", result.Title)
		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&result)
		return
	}

	if r.Body != nil {
		var body todo.Todo
		json.NewDecoder(r.Body).Decode(&body)
		if body.Title == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Print("CreateTodo() title is empty in body")
			return
		}
		log.Printf("recieved todo title: %s", body.Title)
		if body.Title != "" {
			result, err := t.todoStore.CreateTodo(userId, body.Title)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("Error CreateTodo(): %v", err)
				return
			}

			log.Printf("Todo created from json: %v", result.Title)
			w.WriteHeader(http.StatusCreated)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(&result)
			return
		}
	}

	if r.Body == nil && title == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("no title provided by query string and body")
		return
	}

	log.Print("CreateTodo: stange nothing happened")
}

func (t *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	title := r.URL.Query().Get("title")
	completed := r.URL.Query().Get("completed")
	userId, err := auth.GetUserIdFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("auth.GetUserIdFromContext: %v", err)
		return
	}

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
		result, err := t.todoStore.UpdateTodoTitle(userId, id, title)
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
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error parsing id string: %v", err)
			return
		}
		completed, err := strconv.ParseBool(completed)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error parsing completed string to bool: %v", err)
			return
		}
		result, err := t.todoStore.UpdateTodoStatus(userId, id, completed)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error UpdateTodoById(): %v", err)
			return
		}
		json.NewEncoder(w).Encode(&result)
		return
	}

	// json body
	var body todo.Todo
	json.NewDecoder(r.Body).Decode(&body)
	result, err := t.todoStore.UpdateTodoById(userId, body.Id, body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error UpdateTodoById(): %v", err)
		return
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
