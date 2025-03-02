package postgre

import (
	"context"
	"fmt"
	"log"
	"mytodoapp/domain/todo"
	"os"

	"github.com/jackc/pgx/v5"
)

type PostgreTodoStore struct {
	db *pgx.Conn
}

func NewPostgreTodoStore(connString string) *PostgreTodoStore {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return &PostgreTodoStore{db: conn}
}

func (p *PostgreTodoStore) GetTodoByTitle(title string) (todo.Todo, error) {
	var result todo.Todo
	err := p.db.QueryRow(context.Background(), "SELECT title, completed FROM todos WHERE title = $1", title).Scan(&result.Title, &result.Completed)
	if err != nil {
		if err == pgx.ErrNoRows {
			fmt.Fprint(os.Stderr, "No rows")
		}
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return todo.Todo{}, err
	}
	return result, nil
}

func (p *PostgreTodoStore) CreateTodo(title string) (todo.Todo, error) {
	_, err := p.db.Exec(context.Background(), "INSERT INTO todos (title, completed) VALUES ($1, 'false')", title)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exec failed: %v\n", err)
		return todo.Todo{}, err
	}
	return todo.Todo{Title: title, Completed: "false"}, nil
}

func (p *PostgreTodoStore) UpdateTodoTitle(todoToChange string, title string) (todo.Todo, error) {
	_, err := p.db.Exec(context.Background(), "UPDATE todos SET title = $1 WHERE title = $2", title, todoToChange)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exec failed: %v\n", err)
		return todo.Todo{}, err
	}
	return todo.Todo{Title: title, Completed: "false"}, nil
}
