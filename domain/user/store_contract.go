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
		err = sut.RegisterUser(email, password)
		assert.NoError(t, err, "RegisterUser()")
	})
}
