package handler

import (
	"mytodoapp/domain/user"
	"net/http"
)

type UserHandler struct {
	store user.UserStore
}

func NewUserHandler(mux *http.ServeMux, store user.UserStore) {
	handler := &UserHandler{}
	mux.HandleFunc("/login", handler.SignupUser)
}

func (u *UserHandler) SignupUser(w http.ResponseWriter, r *http.Request) {
}
