package todo

import "errors"

var (
	ErrBadRequest    = errors.New("bad request")
	ErrInternalError = errors.New("internal error")
	ErrNotFound      = errors.New("not found")
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
