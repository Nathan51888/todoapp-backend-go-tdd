package httpserver

import (
	"encoding/json"
	"mytodoapp/domain/todo"
	"net/http"
)

type Driver struct {
	BaseURL string
	Client  *http.Client
}

func (d Driver) GetTodoByTitle(title string) (todo.Todo, error) {
	res, err := d.Client.Get(d.BaseURL + "/todo?title=" + title)
	if err != nil {
		return todo.Todo{}, nil
	}
	defer res.Body.Close()

	var result todo.Todo
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return todo.Todo{}, nil
	}
	return result, nil
}
