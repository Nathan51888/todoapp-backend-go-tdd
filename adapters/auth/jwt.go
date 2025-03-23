package auth

import (
	"context"
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
const JWTSecret = "secret"

// TODO: make it an env
const JWTExpirationTime = time.Second * time.Duration(3600*24*3)

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
