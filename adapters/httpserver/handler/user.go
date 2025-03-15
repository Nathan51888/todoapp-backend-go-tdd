package handler

import (
	"encoding/json"
	"log"
	"mytodoapp/domain/user"
	"net/http"
)

type UserHandler struct {
	store user.UserStore
}

func NewUserHandler(mux *http.ServeMux, store user.UserStore) {
	handler := &UserHandler{}
	mux.HandleFunc("/login", handler.LoginUser)
	mux.HandleFunc("/register", handler.RegisterUser)
}

func (u *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
}

func (u *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// get json
	var result any
	json.NewDecoder(r.Body).Decode(&result)
	log.Print(result)

	w.WriteHeader(http.StatusAccepted)

	// check if user exists

	// if not then create user
}
