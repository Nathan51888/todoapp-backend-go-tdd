package handler

import (
	"encoding/json"
	"log"
	"mytodoapp/adapters/auth"
	"mytodoapp/domain/user"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthHandler struct {
	userStore user.UserStore
}

func NewAuthHandler(mux *http.ServeMux, userStore user.UserStore) {
	handler := &AuthHandler{
		userStore: userStore,
	}
	mux.HandleFunc("POST /refresh-token", handler.RefreshToken)
}

func (a *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refreshToken")
	if err != nil {
		log.Printf("error getting refreshToken from cookie: %v", err)
		auth.PermissionDenied(w)
		return
	}

	token, err := auth.ValidateAccessToken(refreshToken.Value)
	if err != nil {
		log.Printf("ValidateAccessToken(): %v", err)
		auth.PermissionDenied(w)
		return
	}
	if !token.Valid {
		log.Printf("invalid token")
		auth.PermissionDenied(w)
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	str := claims["userId"].(string)
	userId, err := uuid.Parse(str)
	if err != nil {
		log.Print("uuid failed to parse userId from token")
		auth.PermissionDenied(w)
		return
	}

	// return new access token
	newAccessToken, err := auth.CreateAccessToken(userId.String())
	if err != nil {
		log.Printf("CreateAccessToken(): %v", err)
		auth.PermissionDenied(w)
		return
	}
	json.NewEncoder(w).Encode(&newAccessToken)
}
