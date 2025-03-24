package postgre

import (
	"context"
	"log"
	"mytodoapp/domain/todo"

	"github.com/google/uuid"
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

func (p *PostgreTodoStore) GetTodoAll(userId uuid.UUID) ([]todo.Todo, error) {
	rows, err := p.db.Query(context.Background(), "SELECT todo_id, title, completed, user_id FROM todos WHERE user_id = $1", userId)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var result []todo.Todo
	for rows.Next() {
		var todo todo.Todo
		if err := rows.Scan(&todo.Id, &todo.Title, &todo.Completed, &todo.UserId); err != nil {
			return nil, err
		}
		result = append(result, todo)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *PostgreTodoStore) GetTodoByTitle(userId uuid.UUID, title string) (todo.Todo, error) {
	var result todo.Todo
	err := p.db.QueryRow(
		context.Background(),
		"SELECT todo_id, title, completed, user_id FROM todos WHERE title = $1",
		title,
	).Scan(&result.Id, &result.Title, &result.Completed, &result.UserId)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println("GetTodoByTitle(): No rows")
		}
		log.Printf("GetTodoByTitle(): QueryRow failed: %v", err)
		return todo.Todo{}, err
	}
	return result, nil
}

func (p *PostgreTodoStore) GetTodoById(userId uuid.UUID, todoId uuid.UUID) (todo.Todo, error) {
	var result todo.Todo
	err := p.db.QueryRow(
		context.Background(),
		"SELECT todo_id, title, completed, user_id FROM todos WHERE todo_id = $1",
		todoId,
	).Scan(&result.Id, &result.Title, &result.Completed, &result.UserId)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println("GetTodoById(): No rows")
		}
		log.Printf("GetTodoById(): QueryRow failed: %v", err)
		return todo.Todo{}, err
	}
	return result, nil
}

func (p *PostgreTodoStore) CreateTodo(userId uuid.UUID, title string) (todo.Todo, error) {
	var result todo.Todo
	err := p.db.QueryRow(
		context.Background(),
		"INSERT INTO todos (todo_id, title, completed, user_id) VALUES (DEFAULT, $1, 'false', $2) RETURNING todo_id",
		title, userId,
	).Scan(&result.Id)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println("CreateTodo(): No rows")
		}
		log.Printf("CreateTodo(): QueryRow failed: %v", err)
		return todo.Todo{}, err
	}
	return todo.Todo{Id: result.Id, Title: title, Completed: false, UserId: userId}, nil
}

func (p *PostgreTodoStore) UpdateTodoTitle(userId uuid.UUID, todoId uuid.UUID, title string) (todo.Todo, error) {
	_, err := p.db.Exec(context.Background(), "UPDATE todos SET title = $1 WHERE todo_id = $2", title, todoId)
	if err != nil {
		log.Printf("UpdateTodoTitle(): Exec failed: %v", err)
		return todo.Todo{}, err
	}
	return todo.Todo{Id: todoId, Title: title, Completed: false, UserId: userId}, nil
}

func (p *PostgreTodoStore) UpdateTodoStatus(userId uuid.UUID, todoId uuid.UUID, completed bool) (todo.Todo, error) {
	_, err := p.db.Exec(
		context.Background(),
		"UPDATE todos SET completed = $1 WHERE todo_id = $2",
		completed, todoId)
	if err != nil {
		log.Printf("UpdateTodoStatus(): Exec failed: %v", err)
		return todo.Todo{}, err
	}
	var result todo.Todo
	err = p.db.QueryRow(
		context.Background(),
		"SELECT todo_id, title, completed, user_id FROM todos WHERE todo_id = $1 AND user_id = $2",
		todoId, userId,
	).Scan(&result.Id, &result.Title, &result.Completed, &result.UserId)
	if err != nil {
		log.Printf("UpdateTodoStatus(): QueryRow failed: %v", err)
		return todo.Todo{}, err
	}
	return result, nil
}

func (p *PostgreTodoStore) UpdateTodoById(userId uuid.UUID, todoId uuid.UUID, todoToChange todo.Todo) (todo.Todo, error) {
	_, err := p.db.Exec(
		context.Background(),
		"UPDATE todos SET title = $1, completed = $2 WHERE todo_id = $3 AND user_id = $4",
		todoToChange.Title, todoToChange.Completed, todoId, userId)
	if err != nil {
		log.Printf("UpdateTodoById(): QueryRow failed: %v", err)
		return todo.Todo{}, err
	}
	return todoToChange, nil
}

func (p *PostgreTodoStore) DeleteTodoById(userId uuid.UUID, todoId uuid.UUID) (todo.Todo, error) {
	var result todo.Todo
	err := p.db.QueryRow(
		context.Background(),
		"DELETE FROM todos WHERE todo_id = $1 AND user_id = $2 RETURNING todo_id, title, completed, user_id",
		todoId, userId,
	).Scan(&result.Id, &result.Title, &result.Completed, &result.UserId)
	if err != nil {
		log.Printf("DeleteTodoById(): QueryRow failed: %v", err)
		return todo.Todo{}, err
	}
	return result, nil
}
