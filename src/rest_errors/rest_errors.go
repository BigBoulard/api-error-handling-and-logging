package rest_errors

import (
	"errors"
	"fmt"
	"net/http"
)

type RestErr interface {
	Status() int     // HTTP status code
	Message() string // Message returned to the client
	Error() string   // Raw Error message
}

type restErr struct {
	ErrStatus  int    `json:"status"`
	ErrMessage string `json:"message"`
	ErrError   error  `json:"error"` // raw error
}

func (e *restErr) Error() string {
	return fmt.Sprintf("message: %s - status: %d - error: %s",
		e.ErrMessage, e.ErrStatus, e.ErrError)
}

func (e *restErr) Message() string {
	return e.ErrMessage
}

func (e *restErr) Status() int {
	return e.ErrStatus
}

// constructors
func NewBadRequestError(message string) RestErr {
	return &restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusBadRequest,
		ErrError:   errors.New("bad_request"),
	}
}

func NewServiceUnavailableError(message string) RestErr {
	return &restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusServiceUnavailable,
		ErrError:   errors.New("service_unavailable"),
	}
}

func NewNotFoundError(message string) RestErr {
	return &restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusNotFound,
		ErrError:   errors.New("not_found"),
	}
}

func NewUnauthorizedError(message string) RestErr {
	return &restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusUnauthorized,
		ErrError:   errors.New("unauthorized"),
	}
}

func NewConflictError(message string) RestErr {
	return &restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusConflict,
		ErrError:   errors.New("conflict"),
	}
}

func NewInternalServerError(message string, err error) RestErr {
	result := &restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusInternalServerError,
		ErrError:   errors.New("internal_server_error"),
	}
	return result
}
