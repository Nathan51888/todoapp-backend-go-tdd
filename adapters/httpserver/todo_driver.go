package httpserver

import (
	"encoding/json"
	"mytodoapp/domain/todo"
	"net/http"

	"github.com/google/uuid"
)

type TodoDriver struct {
	BaseURL string
	Client  *http.Client
	Token   string
}

func (d TodoDriver) GetTodoById(id uuid.UUID) (todo.Todo, error) {
	req, err := http.NewRequest(http.MethodGet, d.BaseURL+"/todo?id="+id.String(), nil)
	if err != nil {
		return todo.Todo{}, err
	}
	req.Header.Add("Authorization", d.Token)
	res, err := d.Client.Do(req)
	if err != nil {
		return todo.Todo{}, err
	}
	defer res.Body.Close()

	var result todo.Todo
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return todo.Todo{}, err
	}
	return result, nil
}

func (d TodoDriver) CreateTodo(title string) (todo.Todo, error) {
	req, err := http.NewRequest(http.MethodPost, d.BaseURL+"/todo?title="+title, nil)
	if err != nil {
		return todo.Todo{}, err
	}
	req.Header.Add("Authorization", d.Token)
	res, err := d.Client.Do(req)
	if err != nil {
		return todo.Todo{}, err
	}
	defer res.Body.Close()

	var result todo.Todo
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return todo.Todo{}, err
	}

	return result, nil
}
