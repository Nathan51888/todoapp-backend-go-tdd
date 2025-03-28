package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"mytodoapp/adapters/auth"
	"mytodoapp/adapters/httpserver/middleware"
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

type UserProfilePayload struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Birthday  string `json:"birthday"`
}

type UserHandler struct {
	store   user.UserStore
	service user.UserService
}

func NewUserHandler(mux *http.ServeMux, store user.UserStore) {
	userService := user.NewUserService(store)
	handler := &UserHandler{store, *userService}
	mux.HandleFunc("POST /login", handler.LoginUser)
	mux.HandleFunc("POST /register", handler.RegisterUser)
	mux.HandleFunc("GET /profile", middleware.WithJWTAuth(handler.GetUser, store))
}

func (u *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var payload LoginUserPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	log.Print("Login payload: ", payload)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokens, err := u.service.LoginUser(payload.Email, payload.Password)
	if err != nil {
		apiError := FromError(err)
		w.WriteHeader(apiError.Status)
		// fmt.Fprint(w, apiError.Message)
		fmt.Fprint(w, "invalid email or password")
		log.Printf("handleLoginUser: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokens.AccessToken, "refreshToken": tokens.RefreshToken})
}

func (u *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	log.Print("Register payload: ", payload)
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
		fmt.Fprint(w, "user email already exists")
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
	userId, err := auth.GetUserIdFromContext(r.Context())
	if err != nil {
		log.Printf("GetUser(): GetUserIdFromContext(): %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := u.store.GetUserById(userId)
	if err != nil {
		log.Printf("error GetUserById: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var userProfilePayload UserProfilePayload
	userProfilePayload.Email = user.Email
	userProfilePayload.FirstName = "first"
	userProfilePayload.LastName = "last"
	userProfilePayload.Birthday = "birthday"

	json.NewEncoder(w).Encode(&userProfilePayload)
}
