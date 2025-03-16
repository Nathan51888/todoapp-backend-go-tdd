package auth_test

import (
	"mytodoapp/adapters/auth"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "password"
	hash, err := auth.HashPassword(password)
	assert.NoError(t, err, "HashPassword()")

	if hash == "" {
		t.Error("expected hash to be not empty")
	}
	if hash == password {
		t.Error("expected hash to be different from original text")
	}
}

func TestComparePassword(t *testing.T) {
	password := "password"
	hash, err := auth.HashPassword(password)
	assert.NoError(t, err, "HashPassword()")

	if !auth.ComparePassword(hash, password) {
		t.Error("expected password to match hash")
	}
	if auth.ComparePassword(hash, "not password") {
		t.Error("expected password to not match hash")
	}
}
