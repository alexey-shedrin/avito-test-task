package httperror

import "net/http"

type HttpError struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	AppError error  `json:"error"`
}

func (e HttpError) Error() string {
	return e.Message
}
func NewInternal(msg string, err error) error {
	return HttpError{Code: http.StatusInternalServerError, Message: msg, AppError: err}
}

func NewBadReq(msg string, err error) error {
	return HttpError{Code: http.StatusBadRequest, Message: msg, AppError: err}
}

func NewUnauthorized(msg string, err error) error {
	return HttpError{Code: http.StatusUnauthorized, Message: msg, AppError: err}
}
