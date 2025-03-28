package user_test

import (
	"errors"
	"mytodoapp/adapters/auth"
	"mytodoapp/adapters/persistence/inmemory"
	"mytodoapp/domain/user"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserLoginHandler(t *testing.T) {
	t.Run("can login with existing user & correct credentials", func(t *testing.T) {
		userPassword, err := auth.HashPassword("password")
		if err != nil {
			t.Fatal(err)
		}
		store := &inmemory.InMemoryUserStore{Users: []user.User{
			{Email: "user@email.com", Password: userPassword},
		}}
		sut := user.NewUserService(store)

		email := "user@email.com"
		password := "password"

		got, err := sut.LoginUser(email, password)
		assert.NoError(t, err)

		if got.AccessToken == "" {
			t.Error("expected token to be not empty")
		}
	})
	t.Run("returns both access and refresh token when logining in", func(t *testing.T) {
		userPassword, err := auth.HashPassword("password")
		if err != nil {
			t.Fatal(err)
		}
		store := &inmemory.InMemoryUserStore{Users: []user.User{
			{Email: "user@email.com", Password: userPassword},
		}}
		sut := user.NewUserService(store)

		email := "user@email.com"
		password := "password"

		got, err := sut.LoginUser(email, password)
		assert.NoError(t, err)
		if got.AccessToken == "" {
			t.Errorf("expected token to be not empty")
		}
		if got.RefreshToken == "" {
			t.Errorf("expected refresh token to not be empty")
		}
	})
	t.Run("returns BadRequest error if wrong password", func(t *testing.T) {
		userPassword, err := auth.HashPassword("correct")
		if err != nil {
			t.Fatal(err)
		}
		store := &inmemory.InMemoryUserStore{Users: []user.User{
			{Email: "user@email.com", Password: userPassword},
		}}
		sut := user.NewUserService(store)
		email := "user@email.com"
		password := "wrong"

		got, err := sut.LoginUser(email, password)
		assert.Error(t, err)
		var serviceError user.Error
		errors.As(err, &serviceError)
		if serviceError.SvcError() != user.ErrBadRequest {
			t.Errorf("expected error type ErrBadRequest: %v", err)
		}
		if got.AccessToken != "" {
			t.Errorf("expected token to be empty")
		}
		if got.RefreshToken != "" {
			t.Errorf("expected refresh token to be empty")
		}
	})
}

func TestUserRegisterHandler(t *testing.T) {
	t.Run("can create user", func(t *testing.T) {
		store := &inmemory.InMemoryUserStore{}
		sut := user.NewUserService(store)
		email := "test@email.com"
		password := "password"

		err := sut.RegisterUser(email, password)
		assert.NoError(t, err)
	})
	t.Run("returns bad request if email is empty", func(t *testing.T) {
		store := &inmemory.InMemoryUserStore{}
		sut := user.NewUserService(store)

		err := sut.RegisterUser("", "")
		assert.Error(t, err)
		var serviceError user.Error
		errors.As(err, &serviceError)
		if serviceError.SvcError() != user.ErrBadRequest {
			t.Errorf("expected error type ErrBadRequest: %v", err)
		}
	})
	t.Run("returns bad request if password is empty", func(t *testing.T) {
		store := &inmemory.InMemoryUserStore{}
		sut := user.NewUserService(store)

		err := sut.RegisterUser("", "")
		assert.Error(t, err)
		var serviceError user.Error
		errors.As(err, &serviceError)
		if serviceError.SvcError() != user.ErrBadRequest {
			t.Errorf("expected error type ErrBadRequest: %v", err)
		}
	})
}

func TestUserProfile(t *testing.T) {
	t.Run("can get user profile with valid token", func(t *testing.T) {
		password, err := auth.HashPassword("password")
		if err != nil {
			t.Fatal(err)
		}
		userId := uuid.New()
		store := &inmemory.InMemoryUserStore{Users: []user.User{
			{Id: userId, Email: "user@email.com", Password: password},
		}}
		sut := user.NewUserService(store)

		got, err := sut.GetUser(userId)
		assert.NoError(t, err)
		want := &user.UserProfile{
			Email:     "user@email.com",
			FirstName: "first",
			LastName:  "last",
			Birthday:  "birthday",
		}
		assert.Equal(t, want, got)
	})
}
