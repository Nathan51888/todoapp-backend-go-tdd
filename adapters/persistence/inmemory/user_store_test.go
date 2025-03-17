package inmemory_test

import (
	"mytodoapp/adapters/persistence/inmemory"
	"mytodoapp/domain/user"
	"testing"
)

func TestInMemoryUserStore(t *testing.T) {
	user.UserStoreContract{NewUserStore: func() (user.UserStore, error) {
		return inmemory.NewInMemoryUserStore()
	}}.Test(t)
}
