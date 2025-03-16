package inmemory

import (
	"errors"
	"mytodoapp/domain/user"

	"github.com/google/uuid"
)

type InMemoryUserStore struct {
	Users []user.User
}

func NewInMemoryUserStore() (*InMemoryUserStore, error) {
	return &InMemoryUserStore{}, nil
}

var test user.UserStore = &InMemoryUserStore{}

func (i *InMemoryUserStore) RegisterUser(email string, password string) (user.User, error) {
	// check if user email exists
	for _, item := range i.Users {
		if item.Email == email {
			return user.User{}, errors.New("user email already exists")
		}
	}

	newUser := user.User{Id: uuid.New(), Email: email, Password: password}
	i.Users = append(i.Users, newUser)
	return newUser, nil
}

func (i *InMemoryUserStore) GetUserByEmail(email string) (user.User, error) {
	for index, item := range i.Users {
		if item.Email == email {
			result := i.Users[index]
			return result, nil
		}
	}
	return user.User{}, errors.New("user email not found")
}

func (i *InMemoryUserStore) GetUserById(id uuid.UUID) (user.User, error) {
	for index, item := range i.Users {
		if item.Id == id {
			result := i.Users[index]
			return result, nil
		}
	}
	return user.User{}, errors.New("user id not found")
}
