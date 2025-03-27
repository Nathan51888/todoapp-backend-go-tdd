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
	JWTExpirationTime        = time.Second * time.Duration(10)
	JWTRefreshExpirationTime = time.Second * time.Duration(3600*24*7)
)

func CreateAccessToken(userId string) (string, error) {
	token := CreateJWT(JWTExpirationTime, userId)
	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token string: %w", err)
	}
	return tokenString, nil
}

func CreateRefreshToken(userId string) (string, error) {
	token := CreateJWT(JWTRefreshExpirationTime, userId)
	tokenString, err := token.SignedString([]byte(JWTRefreshSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token string: %w", err)
	}
	return tokenString, nil
}

func CreateJWT(expiration time.Duration, userId string) *jwt.Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(expiration).Unix(),
	})
	return token
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

	return accessToken, nil
}

func PermissionDenied(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": "permission denied"})
}

func ValidateAccessToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(JWTSecret), nil
	})
	// check for verification errors
	if err != nil {
		return nil, fmt.Errorf("failed to parse/verify token: %w", err)
	}

	// check if token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func ValidateRefreshToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(JWTRefreshSecret), nil
	})
	// check for verification errors
	if err != nil {
		return nil, fmt.Errorf("failed to parse/verify token: %v", err)
	}

	// check if token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
