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

func TestGreeterServer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	var (
		port       = "8080"
		binToBuild = "httpserver"
		baseURL    = fmt.Sprintf("http://localhost:%s", port)
		driver     = httpserver.Driver{BaseURL: baseURL, Client: &http.Client{
			Timeout: 1 * time.Second,
		}}
	)

	adapters.StartDockerServer(t, port, binToBuild)
	specifications.GreetSpecification(t, driver)
}
