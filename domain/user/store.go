package user

type UserStore interface {
	RegisterUser(email string, password string) error
}
