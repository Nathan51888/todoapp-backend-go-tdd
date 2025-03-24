package todo

import "github.com/google/uuid"

type Todo struct {
	Id        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	UserId    uuid.UUID `json:"-"`
}
