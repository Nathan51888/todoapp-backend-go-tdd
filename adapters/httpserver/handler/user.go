package handler

import (
	"encoding/json"
	"log"
	"mytodoapp/adapters/auth"
	"mytodoapp/domain/user"
	"net/http"
)

type RegisterUserPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserHandler struct {
	store user.UserStore
}

func NewUserHandler(mux *http.ServeMux, store user.UserStore) {
	handler := &UserHandler{store}
	mux.HandleFunc("/login", handler.LoginUser)
	mux.HandleFunc("POST /register", handler.RegisterUser)
}

func (u *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
}

func (u *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user RegisterUserPayload
	err := json.NewDecoder(r.Body).Decode(&user)
	log.Print(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = u.store.CreateUser(user.Email, hashedPassword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
