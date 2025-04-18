package handler

import (
	"encoding/json"
	"fmt"
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
	service todo.TodoService
}

func NewTodoHandler(mux *http.ServeMux, todoStore todo.TodoStore, userStore user.UserStore) {
	todoService := todo.NewTodoService(todoStore, userStore)
	handler := &TodoHandler{
		service: *todoService,
	}
	mux.HandleFunc("GET /todo", middleware.WithJWTAuth(handler.handleGetTodo, userStore))
	mux.HandleFunc("POST /todo", middleware.WithJWTAuth(handler.handleCreateTodo, userStore))
	mux.HandleFunc("PUT /todo", middleware.WithJWTAuth(handler.handleUpdateTodo, userStore))
	mux.HandleFunc("DELETE /todo", middleware.WithJWTAuth(handler.handleDeleteTodo, userStore))
}

// GET /todo
func (t *TodoHandler) handleGetTodo(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	id := r.URL.Query().Get("id")

	userId, err := auth.GetUserIdFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("auth.GetUserIdFromContext(): %v", err)
		return
	}

	if title != "" {
		t.GetTodoByTitle(w, r, userId, title)
		return
	}

	if id != "" {
		t.GetTodoById(w, r, userId, id)
		return
	}

	t.GetTodoAll(w, r, userId)
}

func (t *TodoHandler) GetTodoByTitle(w http.ResponseWriter, r *http.Request, userId uuid.UUID, title string) {
	result, err := t.service.GetTodoByTitle(userId, title)
	if err != nil {
		apiError := FromError(err)
		w.WriteHeader(apiError.Status)
		fmt.Fprint(w, apiError.Message)
		log.Printf("handler: GetTodoByTitle(): %v", err)
		return
	}

	json.NewEncoder(w).Encode(&result)
}

func (t *TodoHandler) GetTodoById(w http.ResponseWriter, r *http.Request, userId uuid.UUID, todoIdString string) {
	id, err := uuid.Parse(todoIdString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed converting string to int: %v", err)
		return
	}

	result, err := t.service.GetTodoById(userId, id)
	if err != nil {
		apiError := FromError(err)
		w.WriteHeader(apiError.Status)
		fmt.Fprint(w, apiError.Message)
		log.Printf("handler: GetTodoById(): %v", err)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (t *TodoHandler) GetTodoAll(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	result, err := t.service.GetTodoAll(userId)
	if err != nil {
		apiError := FromError(err)
		w.WriteHeader(apiError.Status)
		fmt.Fprint(w, apiError.Message)
		log.Printf("handler: GetTodoAll(): %v", err)
		return
	}

	json.NewEncoder(w).Encode(&result)
}

// POST /todo
func (t *TodoHandler) handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	userId, err := auth.GetUserIdFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("auth.GetUserIdFromContext: %v", err)
		return
	}

	if title != "" {
		todoToAdd := todo.Todo{Title: title}
		result, err := t.service.CreateTodo(userId, todoToAdd)
		if err != nil {
			apiError := FromError(err)
			w.WriteHeader(apiError.Status)
			fmt.Fprint(w, apiError.Message)
			log.Printf("handleCreateTodo(): %v", err)
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
			todoToAdd := todo.Todo{Title: body.Title}
			result, err := t.service.CreateTodo(userId, todoToAdd)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("CreateTodo(): %v", err)
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

// PUT /todo
func (t *TodoHandler) handleUpdateTodo(w http.ResponseWriter, r *http.Request) {
	todoId := r.URL.Query().Get("id")
	title := r.URL.Query().Get("title")
	completed := r.URL.Query().Get("completed")
	userId, err := auth.GetUserIdFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("auth.GetUserIdFromContext: %v", err)
		return
	}

	if title != "" {
		if todoId == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		todoId, err := uuid.Parse(todoId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed parsing id string: %v", err)
			return
		}
		result, err := t.service.UpdateTodoTitle(userId, todoId, title)
		if err != nil {
			apiError := FromError(err)
			w.WriteHeader(apiError.Status)
			fmt.Fprint(w, apiError.Message)
			log.Printf("handleUpdateTodo(): %v", err)
			return
		}
		json.NewEncoder(w).Encode(&result)
		return
	}

	if completed != "" {
		if todoId == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		todoId, err := uuid.Parse(todoId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed parsing id string: %v", err)
			return
		}
		completed, err := strconv.ParseBool(completed)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed parsing completed string to bool: %v", err)
			return
		}
		result, err := t.service.UpdateTodoStatus(userId, todoId, completed)
		if err != nil {
			apiError := FromError(err)
			w.WriteHeader(apiError.Status)
			fmt.Fprint(w, apiError.Message)
			log.Printf("handleUpdateTodo(): %v", err)
			return
		}
		json.NewEncoder(w).Encode(&result)
		return
	}

	// json body
	var gotTodo todo.Todo
	json.NewDecoder(r.Body).Decode(&gotTodo)
	result, err := t.service.UpdateTodoById(userId, gotTodo.Id, gotTodo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("handleUpdateTodo(): %v", err)
		return
	}
	json.NewEncoder(w).Encode(&result)
}

// DELETE /todo
func (t *TodoHandler) handleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	todoId, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("handleDeleteTodo: %v", err)
		return
	}

	userId, err := auth.GetUserIdFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("auth.GetUserIdFromContext: %v", err)
		return
	}

	result, err := t.service.DeleteTodo(userId, todoId)
	if err != nil {
		apiError := FromError(err)
		w.WriteHeader(apiError.Status)
		fmt.Fprint(w, apiError.Message)
		log.Printf("handleDeleteTodo: %v", err)
		return
	}

	json.NewEncoder(w).Encode(result)
}
