package handler_test

import (
	"bytes"
	"encoding/json"
	"mytodoapp/adapters/httpserver/handler"
	"mytodoapp/adapters/persistence/inmemory"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserHandler(t *testing.T) {
	t.Run("can create user", func(t *testing.T) {
		mux := http.NewServeMux()
		store := &inmemory.InmemoryUserStore{}
		handler.NewUserHandler(mux, store)

		payloadBuff := new(bytes.Buffer)
		json.NewEncoder(payloadBuff).Encode(`[{"email": "test@email.com", "password": "test"}]`)
		req := httptest.NewRequest(http.MethodPost, "/register", payloadBuff)
		res := httptest.NewRecorder()

		mux.ServeHTTP(res, req)

		assert.Equal(t, http.StatusAccepted, res.Code, "status code")
	})
}
