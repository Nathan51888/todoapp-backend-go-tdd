package handler

import (
	"encoding/json"
	"log"
	"mytodoapp/domain/user"
	"net/http"
)

type RegisterUserPayload struct {
	Email    string
	Password string
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
	var payload RegisterUserPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	log.Print(payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = u.store.RegisterUser(payload.Email, payload.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
