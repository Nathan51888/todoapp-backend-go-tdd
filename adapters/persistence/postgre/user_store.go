package postgre

import (
	"context"
	"fmt"
	"log"
	"mytodoapp/domain/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PostgreUserStore struct {
	db *pgx.Conn
}

func NewPostgreUserStore(connString string) (*PostgreUserStore, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	// INFO
	log.Println("UserStore: Connected to database")
	return &PostgreUserStore{db: conn}, nil
}

func (p *PostgreUserStore) CreateUser(email string, password string) (user.User, error) {
	// check if email exists
	var result string
	err := p.db.QueryRow(
		context.Background(),
		"SELECT user_email FROM users WHERE user_email = $1",
		email,
	).Scan(&result)
	// TODO: Might move this to user service
	if result != "" && err != pgx.ErrNoRows {
		return user.User{}, user.ErrUserEmailExists
	}
	if err != nil && err != pgx.ErrNoRows {
		return user.User{}, fmt.Errorf("checking if email exists: %w", err)
	}

	// insert user
	var id uuid.UUID
	err = p.db.QueryRow(
		context.Background(),
		"INSERT INTO users (user_id, user_email, user_password) VALUES (DEFAULT, $1, $2) RETURNING user_id",
		email, password,
	).Scan(&id)
	if err != nil {
		return user.User{}, fmt.Errorf("CreateUser: %v", err)
	}
	return user.User{Id: id, Email: email, Password: password}, nil
}

func (p *PostgreUserStore) GetUserByEmail(email string) (user.User, error) {
	var result user.User
	err := p.db.QueryRow(
		context.Background(),
		"SELECT user_id, user_email, user_password FROM users WHERE user_email = $1",
		email,
	).Scan(&result.Id, &result.Email, &result.Password)
	if err != nil {
		return user.User{}, fmt.Errorf("GetUserByEmail: %v", err)
	}
	return result, nil
}

func (p *PostgreUserStore) GetUserById(id uuid.UUID) (user.User, error) {
	var result user.User
	err := p.db.QueryRow(
		context.Background(),
		"SELECT user_id, user_email, user_password FROM users WHERE user_id = $1",
		id,
	).Scan(&result.Id, &result.Email, &result.Password)
	if err != nil {
		return user.User{}, fmt.Errorf("GetUserById: %v", err)
	}
	return result, nil
}
