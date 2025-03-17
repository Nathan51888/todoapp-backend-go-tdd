package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"mytodoapp/adapters/auth"
	"mytodoapp/domain/user"
	"net/http"

	"github.com/google/uuid"
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
	mux.HandleFunc("GET /profile", auth.WithJWTAuth(handler.GetUser, store))
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

	secret := []byte(auth.JWTSecret)
	token, err := auth.CreateJWT(secret, user.Id.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (u *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	log.Print("payload: ", payload)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// check if user exists
	_, err = u.store.GetUserByEmail(payload.Email)
	if err == nil {
		log.Printf("user with email %s already exists", payload.Email)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = u.store.CreateUser(payload.Email, hashedPassword)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		log.Print("error parsing uuid")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := u.store.GetUserById(userId)
	if err != nil {
		log.Printf("error GetUserById: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&user)
}
