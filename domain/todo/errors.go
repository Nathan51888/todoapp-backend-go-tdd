package todo

import "errors"

var (
	ErrBadRequest     = errors.New("bad request")
	ErrInternalError  = errors.New("internal error")
	ErrNotFound       = errors.New("not found")
	ErrTodoNotFound   = errors.New("todo not found")
	ErrTodoIdEmpty    = errors.New("todo id is empty")
	ErrTodoTitleEmpty = errors.New("title is empty")
	ErrUserIdEmpty    = errors.New("user id is empty")
)

type Error struct {
	appError error
	svcError error
}

func (e Error) AppError() error {
	return e.appError
}

func (e Error) SvcError() error {
	return e.svcError
}

func NewError(svcError, appError error) error {
	return Error{
		svcError: svcError,
		appError: appError,
	}
}

func (e Error) Error() string {
	return errors.Join(e.svcError, e.appError).Error()
}
