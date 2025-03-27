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
	// Read from cookies
	// refreshToken, err := r.Cookie("refreshToken")
	// if err != nil {
	// 	log.Printf("error getting refreshToken from cookie: %v", err)
	// 	auth.PermissionDenied(w)
	// 	return
	// }
	var result map[string]string
	json.NewDecoder(r.Body).Decode(&result)
	refreshToken := result["refreshToken"]
	log.Printf("Refresh token: %v", refreshToken)

	token, err := auth.ValidateRefreshToken(refreshToken)
	if err != nil {
		log.Printf("ValidateRefreshToken(): %v", err)
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
	log.Printf("Created new access token with refresh: %v", newAccessToken)
	json.NewEncoder(w).Encode(map[string]string{"accessToken": newAccessToken})
}
