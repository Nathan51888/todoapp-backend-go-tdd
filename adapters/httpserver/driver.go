package httpserver

import (
	"io"
	"net/http"
)

type Driver struct {
	BaseURL string
	Client  *http.Client
}

func (d Driver) GetTodoByTitle(title string) (string, error) {
	res, err := d.Client.Get(d.BaseURL + "/todo?title=" + title)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	todo, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(todo), nil
}
