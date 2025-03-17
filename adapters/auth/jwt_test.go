package auth_test

import (
	"mytodoapp/adapters/auth"
	"testing"
)

func TestCreateJWT(t *testing.T) {
	secret := []byte(auth.JWTSecret)

	token, err := auth.CreateJWT(secret, "456")
	if err != nil {
		t.Errorf("error creating JWT: %v", err)
	}

	if token == "" {
		t.Error("expected token to be not empty")
	}
}
