package postgre_test

import (
	"context"
	"log"
	"mytodoapp/adapters/persistence/postgre"
	"mytodoapp/domain/todo"
	"testing"
)

func TestPostgreTodoStore(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	pgContainer, err := postgre.CreatePostgresContainer(ctx)
	if err != nil {
		log.Fatalf("failed to create pgContainer: %s", err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})
	userStore, err := postgre.NewPostgreUserStore(pgContainer.ConnectionString)
	if err != nil {
		t.Fatalf("failed to create user store: %v", err)
	}

	todo.TodoStoreContract{NewTodoStore: func() (todo.TodoStore, error) {
		return postgre.NewPostgreTodoStore(pgContainer.ConnectionString)
	}, UserStore: userStore}.Test(t)
}
