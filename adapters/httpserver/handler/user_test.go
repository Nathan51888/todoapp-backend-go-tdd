package handler_test

import (
	"bytes"
	"encoding/json"
	"mytodoapp/adapters/auth"
	"mytodoapp/adapters/httpserver/handler"
	"mytodoapp/adapters/persistence/inmemory"
	"mytodoapp/domain/user"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRegisterHandler(t *testing.T) {
	t.Run("can create user", func(t *testing.T) {
		mux := http.NewServeMux()
		store := &inmemory.InMemoryUserStore{}
		handler.NewUserHandler(mux, store)

		payloadBuff := new(bytes.Buffer)
		user := handler.RegisterUserPayload{"test@email.com", "password"}
		json.NewEncoder(payloadBuff).Encode(&user)
		req := httptest.NewRequest(http.MethodPost, "/register", payloadBuff)
		res := httptest.NewRecorder()

		mux.ServeHTTP(res, req)

		assert.Equal(t, http.StatusCreated, res.Code, "status code")
	})
	t.Run("returns bad request if body is nil", func(t *testing.T) {
		mux := http.NewServeMux()
		store := &inmemory.InMemoryUserStore{}
		handler.NewUserHandler(mux, store)

		req := httptest.NewRequest(http.MethodPost, "/register", nil)
		res := httptest.NewRecorder()

		mux.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})
}

func TestUserLoginHandler(t *testing.T) {
	t.Run("can login with existing user", func(t *testing.T) {
		mux := http.NewServeMux()
		password, err := auth.HashPassword("password")
		if err != nil {
			t.Fatal(err)
		}
		store := &inmemory.InMemoryUserStore{Users: []user.User{
			{Email: "user@email.com", Password: password},
		}}
		handler.NewUserHandler(mux, store)

		payloadBuff := new(bytes.Buffer)
		user := handler.LoginUserPayload{"user@email.com", "password"}
		json.NewEncoder(payloadBuff).Encode(&user)
		req := httptest.NewRequest(http.MethodPost, "/login", payloadBuff)
		res := httptest.NewRecorder()

		mux.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})
}
