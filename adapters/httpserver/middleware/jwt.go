package middleware

import (
	"context"
	"encoding/json"
	"log"
	"mytodoapp/adapters/auth"
	"mytodoapp/domain/user"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func WithJWTAuth(handlerFunc http.HandlerFunc, store user.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := auth.GetAccessTokenFromRequest(r)
		if err != nil {
			log.Printf("GetTokenFromRequest(): %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "missing token"})
			return
		}

		token, err := auth.ValidateAccessToken(tokenString)
		if err != nil {
			log.Printf("failed to validate token ValidateAccessToken(): %v", err)
			auth.PermissionDenied(w)
			return
		}
		if !token.Valid {
			log.Println("invalid token")
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

		u, err := store.GetUserById(userId)
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			auth.PermissionDenied(w)
			return
		}

		// Add the user to the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, auth.UserKey, u.Id)
		r = r.WithContext(ctx)

		// Call the function if the token is valid
		handlerFunc(w, r)
	}
}
