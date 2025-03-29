package user

import (
	"errors"
	"fmt"
	"log"
	"mytodoapp/adapters/auth"

	"github.com/google/uuid"
)

type UserService struct {
	store UserStore
}

type LoginTokens struct {
	AccessToken  string
	RefreshToken string
}

type UserProfile struct {
	Email     string
	FirstName string
	LastName  string
	Birthday  string
}

func NewUserService(store UserStore) *UserService {
	return &UserService{store}
}

func (u *UserService) LoginUser(email string, password string) (*LoginTokens, error) {
	user, err := u.store.GetUserByEmail(email)
	if errors.Is(err, ErrUserEmailExists) {
		return &LoginTokens{}, NewError(ErrBadRequest, err)
	}
	if err != nil {
		return &LoginTokens{}, NewError(ErrInternalError, fmt.Errorf("email not found GetUserByEmail: %w", err))
	}

	if !auth.ComparePassword(user.Password, password) {
		return &LoginTokens{}, NewError(ErrBadRequest, fmt.Errorf("wrong password"))
	}

	accessToken, err := auth.CreateAccessToken(user.Id.String())
	if err != nil {
		return &LoginTokens{}, NewError(ErrInternalError, fmt.Errorf("auth.CreateAccessToken: %w", err))
	}

	refreshToken, err := auth.CreateRefreshToken(user.Id.String())
	if err != nil {
		return &LoginTokens{}, NewError(ErrInternalError, fmt.Errorf("auth.CreateRefreshToken: %w", err))
	}

	return &LoginTokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (u *UserService) RegisterUser(email string, password string) error {
	if email == "" {
		return NewError(ErrBadRequest, fmt.Errorf("email is empty"))
	}
	if password == "" {
		return NewError(ErrBadRequest, fmt.Errorf("password is empty"))
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return NewError(ErrInternalError, fmt.Errorf("auth.HashPassword: %w", err))
	}

	// check if user exists
	_, err = u.store.GetUserByEmail(email)
	if errors.Is(err, ErrUserEmailExists) {
		return NewError(ErrBadRequest, fmt.Errorf("store.GetUserByEmail: %w", err))
	}
	if err == nil {
		log.Printf("user with email %s already exists", email)
		return NewError(ErrInternalError, fmt.Errorf("store.GetUserByEmail: %w", err))
	}

	_, err = u.store.CreateUser(email, hashedPassword)
	if errors.Is(err, ErrUserEmailExists) {
		return NewError(ErrBadRequest, fmt.Errorf("store.CreateUser: %w", err))
	}
	if err != nil {
		return NewError(ErrInternalError, fmt.Errorf("store.CreateUser: %w", err))
	}

	return nil
}

func (u *UserService) GetUser(userId uuid.UUID) (*UserProfile, error) {
	user, err := u.store.GetUserById(userId)
	if err != nil {
		return &UserProfile{}, NewError(ErrInternalError, fmt.Errorf("store.GetUserById: %w", err))
	}

	var userProfile UserProfile
	userProfile.Email = user.Email
	userProfile.FirstName = "first"
	userProfile.LastName = "last"
	userProfile.Birthday = "birthday"

	return &userProfile, nil
}
