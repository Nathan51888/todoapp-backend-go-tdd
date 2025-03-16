package handler

import (
	"encoding/json"
	"fmt"
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
	mux.HandleFunc("POST /login", handler.LoginUser)
	mux.HandleFunc("POST /register", handler.RegisterUser)
}

func (u *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var payload LoginUserPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	log.Print("payload: ", payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := u.store.GetUserByEmail(payload.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("GetUserByEmail(): not found, invalid email or password")
		fmt.Fprint(w, "not found, invalid email or password")
		return
	}

	if !auth.ComparePassword(user.Password, payload.Password) {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("ComparePassword(): not found, invalid email or password")
		fmt.Fprint(w, "not found, invalid email or password")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"secret": ""})
}

func (u *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	log.Print("payload: ", payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = u.store.CreateUser(payload.Email, hashedPassword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
