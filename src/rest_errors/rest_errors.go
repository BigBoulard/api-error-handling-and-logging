package rest_errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type RestErr interface {
	Status() int     // HTTP status code
	Message() string // Message returned to the client
	Error() string   // Raw Error message
	Causes() string  // pkg - method | pkg - method for logging purposes
	AddCause(string)
}

type restErr struct {
	ErrStatus  int      `json:"status"`
	ErrMessage string   `json:"message"`
	ErrError   string   `json:"error"`
	ErrCauses  []string `json:"causes"`
}

func (e *restErr) Error() string {
	return fmt.Sprintf("message: %s - status: %d - error: %s - causes: %v",
		e.ErrMessage, e.ErrStatus, e.ErrError, e.ErrCauses)
}

func (e *restErr) Message() string {
	return e.ErrMessage
}

func (e *restErr) Status() int {
	return e.ErrStatus
}

func (e *restErr) Causes() string {
	return strings.Join(e.ErrCauses, " | ")
}

func (e *restErr) AddCause(cause string) {
	e.ErrCauses = prependString(e.ErrCauses, cause)
}

func NewRestError(message string, status int, err string, causes []string) RestErr {
	return &restErr{
		ErrMessage: message,
		ErrStatus:  status,
		ErrError:   err,
		ErrCauses:  causes,
	}
}

func NewRestErrorFromBytes(bytes []byte) (RestErr, error) {
	var apiErr restErr
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("invalid json")
	}
	return &apiErr, nil
}

func NewBadRequestError(message string) RestErr {
	return &restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusBadRequest,
		ErrError:   "bad_request",
	}
}

func NewServiceUnavailableError(message string) RestErr {
	return &restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusServiceUnavailable,
		ErrError:   "service_unavailable",
	}
}

func NewNotFoundError(message string) RestErr {
	return &restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusNotFound,
		ErrError:   "not_found",
	}
}

func NewGoneError(message string) RestErr {
	return &restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusGone,
		ErrError:   "gone",
	}
}

func NewUnauthorizedError(message string) RestErr {
	return &restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusUnauthorized,
		ErrError:   "unauthorized",
	}
}

func NewConflictError(message string) RestErr {
	return &restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusConflict,
		ErrError:   "conflict",
	}
}

func NewUnprocessableEntityError(message string) RestErr {
	return &restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusUnprocessableEntity,
		ErrError:   "unprocessable_entity",
	}
}

func NewInternalServerError(message string, err error) RestErr {
	result := &restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusInternalServerError,
		ErrError:   "internal_server_error",
	}
	if err != nil {
		result.ErrCauses = prependString(result.ErrCauses, err.Error())
	}
	return result
}

func prependString(x []string, y string) []string {
	x = append(x, "")
	copy(x[1:], x)
	x[0] = y
	return x
}