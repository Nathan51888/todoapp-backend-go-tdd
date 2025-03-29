package httpserver

import (
	"errors"
	"mytodoapp/domain/user"
	"net/http"
)

type APIError struct {
	Status  int
	Message string
}

func FromError(err error) APIError {
	var apiError APIError
	var serviceError user.Error
	if errors.As(err, &serviceError) {
		apiError.Message = serviceError.SvcError().Error()

		svcErr := serviceError.SvcError()
		switch svcErr {
		case user.ErrBadRequest:
			apiError.Status = http.StatusBadRequest
		case user.ErrInternalError:
			apiError.Status = http.StatusInternalServerError
		case user.ErrNotFound:
			apiError.Status = http.StatusNotFound
		}
	}
	return apiError
}
