package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const UserKey contextKey = "userId"

// TODO: make it an env
const (
	JWTSecret        = "access"
	JWTRefreshSecret = "refresh"
)

// TODO: make it an env
const (
	JWTExpirationTime        = time.Second * time.Duration(60*5)
	JWTRefreshExpirationTime = time.Second * time.Duration(3600*24*7)
)

func CreateAccessToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    userId,
		"expiredAt": time.Now().Add(JWTExpirationTime).Unix(),
	})

	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func CreateRefreshToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    userId,
		"expiredAt": time.Now().Add(JWTRefreshExpirationTime).Unix(),
	})

	tokenString, err := token.SignedString([]byte(JWTRefreshSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetUserIdFromContext(ctx context.Context) (uuid.UUID, error) {
	userId, ok := ctx.Value(UserKey).(uuid.UUID)
	if !ok {
		log.Print("Couldn't get userId from context")
		return uuid.UUID{}, errors.New("couldn't get userId from context")
	}

	return userId, nil
}

func GetAccessTokenFromRequest(r *http.Request) (string, error) {
	accessToken := r.Header.Get("Authorization")

	if accessToken == "" {
		log.Print("No accessToken was recieved in header")
		return "", fmt.Errorf("no access token in header")
	}

	log.Printf("Access token recieved: %v", accessToken)
	return accessToken, nil
}

func PermissionDenied(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": "permission denied"})
}

func ValidateAccessToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(JWTSecret), nil
	})
}

func ValidateRefreshToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(JWTRefreshSecret), nil
	})
}
