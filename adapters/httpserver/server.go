package httpserver

import (
	"encoding/json"
	"log"
	"mytodoapp/domain/todo"
	"net/http"
	"strconv"
)

type TodoService interface {
	GetTodoByTitle(title string) (todo.Todo, error)
}

type TodoServer struct {
	store todo.TodoStore
	http.Handler
}

func NewTodoServer(store todo.TodoStore) *TodoServer {
	server := new(TodoServer)

	server.store = store

	mux := http.NewServeMux()
	mux.HandleFunc("GET /todo", server.GetTodo)
	mux.HandleFunc("POST /todo", server.CreateTodo)
	mux.HandleFunc("PUT /todo", server.UpdateTodo)

	server.Handler = mux

	return server
}

func (t *TodoServer) GetTodo(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	id := r.URL.Query().Get("id")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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

func (t *TodoServer) GetTodoByTitle(w http.ResponseWriter, r *http.Request, title string) {
	result, err := t.store.GetTodoByTitle(title)
	if err != nil {
		log.Printf("Error GetTodoByTitle(): %v", err)
	}

	json.NewEncoder(w).Encode(&result)
}

func (t *TodoServer) GetTodoById(w http.ResponseWriter, r *http.Request, idString string) {
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Printf("Error converting string to int: %v", err)
	}

	result, err := t.store.GetTodoById(id)
	if err != nil {
		log.Printf("Error GetTodoById(): %v", err)
	}

	json.NewEncoder(w).Encode(result)
}

func (t *TodoServer) GetTodoAll(w http.ResponseWriter, r *http.Request) {
	result, err := t.store.GetTodoAll()
	if err != nil {
		log.Printf("Error GetTodoAll(): %v", err)
	}

	json.NewEncoder(w).Encode(&result)
}

func (t *TodoServer) CreateTodo(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")

	if title != "" {
		result, err := t.store.CreateTodo(title)
		if err != nil {
			log.Printf("Error CreateTodo(): %v", err)
		}
		json.NewEncoder(w).Encode(result)
		return
	}

	var body todo.Todo
	json.NewDecoder(r.Body).Decode(&body)

	result, err := t.store.CreateTodo(body.Title)
	if err != nil {
		log.Printf("Error CreateTodo(): %v", err)
	}

	json.NewEncoder(w).Encode(result)
}

func (t *TodoServer) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	var body todo.Todo
	json.NewDecoder(r.Body).Decode(&body)

	result, err := t.store.UpdateTodoById(body.Id, body)
	if err != nil {
		log.Printf("Error UpdateTodoById(): %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(&result)
}
