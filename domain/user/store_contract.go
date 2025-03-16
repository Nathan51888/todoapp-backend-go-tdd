package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserStoreContract struct {
	NewUserStore func() (UserStore, error)
}

func (c UserStoreContract) Test(t *testing.T) {
	t.Run("can register user", func(t *testing.T) {
		sut, err := c.NewUserStore()
		if err != nil {
			t.Fatalf("Error creating UserStore: %v", err)
		}

		email := "email@email.com"
		password := "password"
		registeredUser, err := sut.RegisterUser(email, password)
		assert.NoError(t, err, "RegisterUser()")
		assert.Equal(t, email, registeredUser.Email)
		assert.Equal(t, password, registeredUser.Password)

		foundUser, err := sut.GetUserById(registeredUser.Id)
		assert.NoError(t, err, "GetUserById()")
		assert.Equal(t, registeredUser, foundUser, "GetUserById()")
	})
	t.Run("can get registered user by email", func(t *testing.T) {
		sut, err := c.NewUserStore()
		if err != nil {
			t.Fatalf("Error creating UserStore: %v", err)
		}

		email := "email@email.com"
		password := "password"
		registeredUser, err := sut.RegisterUser(email, password)
		assert.NoError(t, err, "RegisterUser()")
		assert.Equal(t, email, registeredUser.Email)
		assert.Equal(t, password, registeredUser.Password)

		foundUser, err := sut.GetUserByEmail(email)
		assert.NoError(t, err, "GetUserByEmail()")
		assert.Equal(t, registeredUser, foundUser, "GetUserByEmail()")
	})
	t.Run("cannot register user if email exists", func(t *testing.T) {
		sut, err := c.NewUserStore()
		if err != nil {
			t.Fatalf("Error creating UserStore: %v", err)
		}

		email := "existing@email.com"
		password := "password"
		_, err = sut.RegisterUser(email, password)
		assert.NoError(t, err, "RegisterUser()")
		_, err = sut.RegisterUser(email, password)
		assert.Error(t, err, "RegisterUser()")
	})
}
