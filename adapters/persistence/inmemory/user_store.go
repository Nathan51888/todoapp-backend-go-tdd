package inmemory

import "mytodoapp/domain/user"

type User struct {
	email    string
	password string
}

type InMemoryUserStore struct {
	Users []User
}

func NewInMemoryUserStore() (*InMemoryUserStore, error) {
	return &InMemoryUserStore{}, nil
}

var test user.UserStore = &InMemoryUserStore{}

func (i *InMemoryUserStore) RegisterUser(email string, password string) error {
	i.Users = append(i.Users, User{email, password})
	return nil
}
