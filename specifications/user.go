package specifications

import (
	"mytodoapp/adapters/httpserver/handler"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserDriver interface {
	RegisterUser(email string, password string) error
	LoginUser(email string, password string) (token string, error error)
	GetUserProfile(token string) (handler.UserProfilePayload, error)
}

func UserSpecification(t testing.TB, driver UserDriver) {
	email := "test@email.com"
	password := "test"
	err := driver.RegisterUser(email, password)
	assert.NoError(t, err, "driver.RegisterUser()")

	token, err := driver.LoginUser(email, password)
	assert.NoError(t, err, "driver.LoginUser()")
	if token == "" {
		t.Error("error driver.LoginUser(): expected token to be not empty")
	}

	want := handler.UserProfilePayload{
		Email:     "test@email.com",
		FirstName: "first",
		LastName:  "last",
		Birthday:  "birthday",
	}
	profile, err := driver.GetUserProfile(token)
	assert.NoError(t, err, "driver.GetUserProfile()")
	assert.Equal(t, want, profile)
}
