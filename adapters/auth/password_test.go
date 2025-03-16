package auth_test

import (
	"mytodoapp/adapters/auth"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "password"
	hashedPassword, err := auth.HashPassword(password)
	assert.NoError(t, err, "HashPassword()")
	isCorrect := auth.ComparePassword(hashedPassword, password)
	assert.Equal(t, true, isCorrect)
}
