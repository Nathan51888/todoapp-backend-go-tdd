package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mytodoapp/domain/user"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const UserKey contextKey = "userId"

// TODO: make it an env
const JWTSecret = "secret"

// TODO: make it an env
const JWTExpirationTime = time.Second * time.Duration(3600*24*3)

// TODO: make it a middleware
func WithJWTAuth(handlerFunc http.HandlerFunc, store user.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := GetTokenFromRequest(r)

		token, err := validateJWT(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["userId"].(string)
		userId, err := uuid.Parse(str)
		if err != nil {
			log.Print("uuid failed to parse userId from token")
			permissionDenied(w)
			return
		}

		u, err := store.GetUserById(userId)
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			permissionDenied(w)
			return
		}

		// Add the user to the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.Id)
		r = r.WithContext(ctx)

		// Call the function if the token is valid
		handlerFunc(w, r)
	}
}

func CreateJWT(secret []byte, userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    userId,
		"expiredAt": time.Now().Add(JWTExpirationTime).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret"), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(map[string]string{"error": "permission denied"})
}

func GetUserIdFromContext(ctx context.Context) (uuid.UUID, error) {
	userId, ok := ctx.Value(UserKey).(uuid.UUID)
	if !ok {
		log.Print("Couldn't get userId from context")
		return uuid.UUID{}, errors.New("couldn't get userId from context")
	}

	return userId, nil
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	log.Printf("Cookies recieved: %v", r.Cookies())

	if tokenAuth != "" {
		log.Printf("Header recieved: %v", tokenAuth)
		return tokenAuth
	}

	log.Print("No header was recieved")
	return ""
}
