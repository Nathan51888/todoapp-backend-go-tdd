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
	service user.UserService
}

func NewUserHandler(mux *http.ServeMux, userStore user.UserStore) {
	userService := user.NewUserService(userStore)
	handler := &UserHandler{
		service: *userService,
	}
	mux.HandleFunc("POST /login", handler.handleLoginUser)
	mux.HandleFunc("POST /register", handler.handleRegisterUser)
	mux.HandleFunc("GET /profile", middleware.WithJWTAuth(handler.handleGetUser, userStore))
}

func (u *UserHandler) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	var payload LoginUserPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	// FIXME: don't exapose password
	log.Print("Login payload: ", payload)
	if err != nil {
		log.Printf("handleLoginUser: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokens, err := u.service.LoginUser(payload.Email, payload.Password)
	if err != nil {
		apiError := FromError(err)
		w.WriteHeader(apiError.Status)
		fmt.Fprint(w, "invalid email or password")
		if apiError.Status == http.StatusInternalServerError {
			fmt.Fprint(w, apiError.Message)
		}
		log.Printf("handleLoginUser: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokens.AccessToken, "refreshToken": tokens.RefreshToken})
}

func (u *UserHandler) handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	// FIXME: don't exapose password
	log.Print("Register payload: ", payload)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = u.service.RegisterUser(payload.Email, payload.Password)
	if err != nil {
		apiError := FromError(err)
		w.WriteHeader(apiError.Status)
		fmt.Fprint(w, apiError.Message)
		log.Printf("handleRegisterUser: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (u *UserHandler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserIdFromContext(r.Context())
	if err != nil {
		log.Printf("GetUser(): GetUserIdFromContext(): %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := u.service.GetUser(userId)
	if err != nil {
		apiError := FromError(err)
		w.WriteHeader(apiError.Status)
		fmt.Fprint(w, apiError.Message)
		log.Printf("handleGetUserById: %v", err)
		return
	}

	var userProfilePayload UserProfilePayload
	userProfilePayload.Email = user.Email
	userProfilePayload.FirstName = "first"
	userProfilePayload.LastName = "last"
	userProfilePayload.Birthday = "birthday"

	json.NewEncoder(w).Encode(&userProfilePayload)
}
