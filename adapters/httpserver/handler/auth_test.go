package handler_test

import (
	"bytes"
	"encoding/json"
	"log"
	"mytodoapp/adapters/auth"
	"mytodoapp/adapters/httpserver/handler"
	"mytodoapp/adapters/persistence/inmemory"
	"mytodoapp/domain/user"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler(t *testing.T) {
	userId := uuid.New()
	userStore := &inmemory.InMemoryUserStore{Users: []user.User{
		{Id: userId},
	}}
	handler := createAuthHandler(userStore)
	refreshToken, err := auth.CreateRefreshToken(userId.String())
	if err != nil {
		t.Fatalf("failed creating refresh token: %v", err)
	}

	t.Run("returns access token with valid refresh token", func(t *testing.T) {
		payloadBuff := new(bytes.Buffer)
		payload := map[string]string{"refreshToken": refreshToken}
		json.NewEncoder(payloadBuff).Encode(&payload)
		req := httptest.NewRequest(http.MethodPost, "/refresh-token", payloadBuff)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		var got map[string]string
		json.NewDecoder(res.Body).Decode(&got)
		log.Println("Body: ", got)
		token, err := auth.ValidateAccessToken(got["accessToken"])
		if err != nil {
			t.Fatalf("error ValidateAccessToken(): %v", err)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userId"].(string)
		id, err := uuid.Parse(str)
		if err != nil {
			t.Errorf("error parsing uuid string: %v", err)
		}
		assert.Equal(t, userId, id)
	})

	t.Run("returns 401 if refresh token is expired", func(t *testing.T) {
		token := auth.CreateJWT(time.Duration(-50)*time.Second, userId.String())
		tokenString, err := token.SignedString([]byte(auth.JWTRefreshSecret))
		if err != nil {
			t.Fatalf("error signing token: %v", err)
		}

		tokenCookie := http.Cookie{
			Name:     "refreshToken",
			Value:    tokenString,
			Expires:  time.Now().Add(time.Hour * time.Duration(24*7)),
			Secure:   true,
			HttpOnly: true,
		}
		req := httptest.NewRequest(http.MethodPost, "/refresh-token", nil)
		req.AddCookie(&tokenCookie)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
		log.Println(res.Body)
	})
}

func createAuthHandler(userStore user.UserStore) *http.ServeMux {
	mux := http.NewServeMux()
	handler.NewAuthHandler(mux, userStore)
	return mux
}
