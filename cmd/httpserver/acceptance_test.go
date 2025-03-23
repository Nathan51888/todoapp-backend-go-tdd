package main_test

import (
	"context"
	"fmt"
	"mytodoapp/adapters/httpserver"
	"mytodoapp/specifications"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
)

func TestTodoServer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	// TODO:Use dynamic url for containers
	compose, err := tc.NewDockerCompose("../../compose.yaml")
	require.NoError(t, err, "NewDockerComposeAPI()")

	t.Cleanup(func() {
		require.NoError(t, compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal), "compose.Down()")
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	require.NoError(t, compose.Up(ctx, tc.Wait(true)), "compose.Up()")

	var (
		backendPort = "8080"
		baseURL     = fmt.Sprintf("http://localhost:%s", backendPort)
	)

	userDriver := httpserver.UserDriver{BaseURL: baseURL, Client: &http.Client{
		Timeout: 1 * time.Second,
	}}
	specifications.UserSpecification(t, &userDriver)

	todoDriver := httpserver.TodoDriver{BaseURL: baseURL, Client: &http.Client{
		Timeout: 1 * time.Second,
	}, Token: userDriver.Token}
	specifications.TodoSpecification(t, &todoDriver)
}
