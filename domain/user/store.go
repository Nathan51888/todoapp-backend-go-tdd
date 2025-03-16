package user

import "github.com/google/uuid"

type UserStore interface {
	RegisterUser(email string, password string) (User, error)
	GetUserByEmail(email string) (User, error)
	GetUserById(id uuid.UUID) (User, error)
}
