package postgre

import (
	"context"
	"log"
	"mytodoapp/domain/todo"

	"github.com/jackc/pgx/v5"
)

type PostgreTodoStore struct {
	db *pgx.Conn
}

func NewPostgreTodoStore(connString string) (*PostgreTodoStore, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Printf("Unable to connect to database: %v", err)
		return nil, err
	}
	log.Println("Connected to database")
	return &PostgreTodoStore{db: conn}, nil
}

func (p *PostgreTodoStore) GetTodoByTitle(title string) (todo.Todo, error) {
	var result todo.Todo
	err := p.db.QueryRow(context.Background(), "SELECT title, completed FROM todos WHERE title = $1", title).Scan(&result.Title, &result.Completed)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println("GetTodoByTitle(): No rows")
		}
		log.Printf("QueryRow failed: %v", err)
		return todo.Todo{}, err
	}
	return result, nil
}

func (p *PostgreTodoStore) CreateTodo(title string) (todo.Todo, error) {
	_, err := p.db.Exec(context.Background(), "INSERT INTO todos (title, completed) VALUES ($1, 'false')", title)
	if err != nil {
		log.Printf("Exec failed: %v", err)
		return todo.Todo{}, err
	}
	return todo.Todo{Title: title, Completed: "false"}, nil
}

func (p *PostgreTodoStore) UpdateTodoTitle(todoToChange string, title string) (todo.Todo, error) {
	_, err := p.db.Exec(context.Background(), "UPDATE todos SET title = $1 WHERE title = $2", title, todoToChange)
	if err != nil {
		log.Printf("Exec failed: %v", err)
		return todo.Todo{}, err
	}
	return todo.Todo{Title: title, Completed: "false"}, nil
}

func (p *PostgreTodoStore) UpdateTodoStatus(todoToChange string, completed string) (todo.Todo, error) {
	_, err := p.db.Exec(context.Background(), "UPDATE todos SET completed = $1 WHERE title = $2", completed, todoToChange)
	if err != nil {
		log.Printf("Exec failed: %v", err)
		return todo.Todo{}, err
	}
	return todo.Todo{Title: todoToChange, Completed: completed}, nil
}
