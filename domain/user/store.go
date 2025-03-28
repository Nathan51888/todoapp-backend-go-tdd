package user

import (
	"errors"

	"github.com/google/uuid"
)

var ErrUserEmailExists = errors.New("user email already exists")

type UserStore interface {
	CreateUser(email string, password string) (User, error)
	GetUserByEmail(email string) (User, error)
	GetUserById(id uuid.UUID) (User, error)
}
