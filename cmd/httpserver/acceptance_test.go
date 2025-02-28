package main_test

import (
	"fmt"
	"mytodoapp/adapters"
	"mytodoapp/adapters/httpserver"
	"mytodoapp/specifications"
	"net/http"
	"testing"
	"time"
)

func TestTodoServer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	var (
		port       = "8080"
		baseURL    = fmt.Sprintf("http://localhost:%s", port)
		binToBuild = "httpserver"
		driver     = httpserver.Driver{BaseURL: baseURL, Client: &http.Client{
			Timeout: 1 * time.Second,
		}}
	)

	adapters.StartDockerServer(t, port, binToBuild)
	specifications.TodoSpecification(t, driver)
}
