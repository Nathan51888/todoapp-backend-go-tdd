package postgre_test

import (
	"context"
	"log"
	"mytodoapp/adapters/persistence/postgre"
	"mytodoapp/domain/user"
	"testing"
)

func TestPostgresUserStore(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	pgContainer, err := postgre.CreatePostgresContainer(ctx)
	if err != nil {
		log.Fatalf("failed to create pgContainer: %v", err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %v", err)
		}
	})

	user.UserStoreContract{NewUserStore: func() (user.UserStore, error) {
		return postgre.NewPostgreUserStore(pgContainer.ConnectionString)
	}}.Test(t)
}
