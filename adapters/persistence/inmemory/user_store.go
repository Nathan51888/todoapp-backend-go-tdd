package inmemory

import "mytodoapp/domain/user"

type InmemoryUserStore struct {
	Users []string
}

var test user.UserStore = &InmemoryUserStore{}

func (i *InmemoryUserStore) GetUser() {
}
