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

func (p *PostgreTodoStore) GetTodoAll() ([]todo.Todo, error) {
	rows, err := p.db.Query(context.Background(), "SELECT id, title, completed FROM todos")
	if err != nil {
		log.Printf("Query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var result []todo.Todo
	for rows.Next() {
		var todo todo.Todo
		if err := rows.Scan(&todo.Id, &todo.Title, &todo.Completed); err != nil {
			return nil, err
		}
		result = append(result, todo)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *PostgreTodoStore) GetTodoByTitle(title string) (todo.Todo, error) {
	var result todo.Todo
	err := p.db.QueryRow(context.Background(), "SELECT * FROM todos WHERE title = $1", title).Scan(&result.Id, &result.Title, &result.Completed)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println("GetTodoByTitle(): No rows")
		}
		log.Printf("GetTodoByTitle(): QueryRow failed: %v", err)
		return todo.Todo{}, err
	}
	return result, nil
}

func (p *PostgreTodoStore) GetTodoById(todoId int) (todo.Todo, error) {
	var result todo.Todo
	err := p.db.QueryRow(context.Background(), "SELECT * FROM todos WHERE id = $1", todoId).Scan(&result.Id, &result.Title, &result.Completed)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println("GetTodoById(): No rows")
		}
		log.Printf("GetTodoById(): QueryRow failed: %v", err)
		return todo.Todo{}, err
	}
	return result, nil
}

func (p *PostgreTodoStore) CreateTodo(title string) (todo.Todo, error) {
	var result todo.Todo
	err := p.db.QueryRow(context.Background(), "INSERT INTO todos (id, title, completed) VALUES (DEFAULT, $1, 'false') RETURNING id", title).Scan(&result.Id)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println("CreateTodo(): No rows")
		}
		log.Printf("CreateTodo(): QueryRow failed: %v", err)
		return todo.Todo{}, err
	}
	return todo.Todo{Id: result.Id, Title: title, Completed: false}, nil
}

func (p *PostgreTodoStore) UpdateTodoTitle(todoId int, title string) (todo.Todo, error) {
	_, err := p.db.Exec(context.Background(), "UPDATE todos SET title = $1 WHERE id = $2", title, todoId)
	if err != nil {
		log.Printf("UpdateTodoTitle(): Exec failed: %v", err)
		return todo.Todo{}, err
	}
	return todo.Todo{Id: todoId, Title: title, Completed: false}, nil
}

func (p *PostgreTodoStore) UpdateTodoStatus(todoId int, completed bool) (todo.Todo, error) {
	_, err := p.db.Exec(context.Background(), "UPDATE todos SET completed = $1 WHERE id = $2", completed, todoId)
	if err != nil {
		log.Printf("UpdateTodoStatus(): Exec failed: %v", err)
		return todo.Todo{}, err
	}
	var result todo.Todo
	err = p.db.QueryRow(context.Background(), "SELECT * FROM todos WHERE id = $1", todoId).Scan(&result.Id, &result.Title, &result.Completed)
	if err != nil {
		log.Printf("UpdateTodoStatus(): QueryRow failed: %v", err)
		return todo.Todo{}, err
	}
	return result, nil
}

func (p *PostgreTodoStore) UpdateTodoById(todoId int, todoToChange todo.Todo) (todo.Todo, error) {
	_, err := p.db.Exec(context.Background(), "UPDATE todos SET title = $1, completed = $2 WHERE id = $3", todoToChange.Title, todoToChange.Completed, todoId)
	if err != nil {
		log.Printf("UpdateTodoById(): QueryRow failed: %v", err)
		return todo.Todo{}, err
	}
	return todoToChange, nil
}

func (p *PostgreTodoStore) DeleteTodoById(todoId int) (todo.Todo, error) {
	var result todo.Todo
	err := p.db.QueryRow(context.Background(), "DELETE FROM todos WHERE id = $1 RETURNING id, title, completed", todoId).Scan(&result.Id, &result.Title, &result.Completed)
	if err != nil {
		log.Printf("DeleteTodoById(): QueryRow failed: %v", err)
		return todo.Todo{}, err
	}
	return result, nil
}
